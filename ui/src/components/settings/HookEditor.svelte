<script>
    import { onMount } from "svelte";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";
    import PageWrapper from "@/components/base/PageWrapper.svelte";

    let files = [];
    let selectedFile = null;
    let editor;
    let editorContainer;
    const baseUrl = import.meta.env.PB_BACKEND_URL;

    async function loadDtsFile(url) {
        const response = await fetch(url);
        const dtsContent = await response.text();
        window.monaco.languages.typescript.javascriptDefaults.addExtraLib(dtsContent, url);
    }

    async function loadMonacoEditor() {
        if (window.monaco) return;

        const script = document.createElement("script");
        script.src = "./libs/monaco-editor/min/vs/loader.js";
        script.async = true;
        document.body.appendChild(script);

        await new Promise((resolve) => (script.onload = resolve));

        require.config({ paths: { vs: "./libs/monaco-editor/min/vs" } });

        await new Promise((resolve) => {
            require(["vs/editor/editor.main"], resolve);
        });
    }

    onMount(async () => {
        await loadMonacoEditor();

        loadFiles();

        editor = window.monaco.editor.create(editorContainer, {
            value: "",
            language: "javascript",
            theme: "vs",
        });

        await loadDtsFile(baseUrl + "types.d.ts");

        return () => {
            editor.dispose();
        };
    });

    async function authenticatedFetch(url, options = {}) {
        const authData = localStorage.getItem("__pb_superuser_auth__");
        let token = "";
        if (authData) {
            try {
                const parsedData = JSON.parse(authData);
                token = parsedData.token;
            } catch (e) {
                goto("/login");
            }
        }

        const headers = new Headers(options.headers || {});
        headers.set("Authorization", token);

        return fetch(url, {
            ...options,
            headers: headers,
        });
    }

    async function loadFiles() {
        const response = await authenticatedFetch(baseUrl + "api/hooks");
        const result = await response.json();

        if (result) {
            files = result;
        }
    }

    async function selectFile(filename) {
        selectedFile = filename;
        const response = await authenticatedFetch(baseUrl + `api/hooks/${filename}`);
        const content = await response.json();
        editor.setValue(content);
    }

    function encryptString(inputString) {
        // 将输入字符串转换为Base64
        const base64String = btoa(encodeURIComponent(inputString));

        // 对Base64字符串进行简单加密
        let encrypted = "";
        for (let i = 0; i < base64String.length; i++) {
            const char = base64String[i];
            if (char.match(/[a-zA-Z]/)) {
                // 对字母进行凯撒密码加密（偏移量为3）
                const code = base64String.charCodeAt(i);
                if (char === char.toUpperCase()) {
                    encrypted += String.fromCharCode(((code - 65 + 3) % 26) + 65);
                } else {
                    encrypted += String.fromCharCode(((code - 97 + 3) % 26) + 97);
                }
            } else if (char.match(/[0-9]/)) {
                // 对数字进行简单替换
                encrypted += String.fromCharCode(((char.charCodeAt(0) - 48 + 5) % 10) + 48);
            } else {
                // 保持其他字符不变
                encrypted += char;
            }
        }

        return encrypted;
    }

    async function saveFile() {
        if (!selectedFile) return;

        const content = editor.getValue();
        const base64Content = encryptString(content);

        await authenticatedFetch(baseUrl + `api/hooks/${selectedFile}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ content: base64Content }),
        });

        addSuccessToast("文件保存成功");
    }

    async function createNewFile() {
        const filename = prompt("输入文件名:");
        if (!filename) return;

        const response = await authenticatedFetch(baseUrl + `api/hooks/`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ filename: filename }),
        });

        if (response.ok) {
            addSuccessToast("创建文件成功");
            await loadFiles(); // Refresh the file list
            selectFile(filename); // Select the new file
        } else {
            addErrorToast("创建失败：" + (await response.json()).message);
        }
    }

    async function renameFile() {
        if (!selectedFile) return;

        const newFilename = prompt("新文件名:", selectedFile);
        if (!newFilename || newFilename === selectedFile) return;

        const response = await authenticatedFetch(baseUrl + `api/hooks/`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ oldname: selectedFile, newname: newFilename }),
        });

        if (response.ok) {
            addSuccessToast("重命名成功");
            await loadFiles();
            selectFile(newFilename);
        } else {
            addErrorToast("重命名失败：" + (await response.json()).message);
        }
    }

    async function deleteFile() {
        if (!selectedFile) return;

        confirm(`确认要删除 ${selectedFile} 吗？`, async () => {
            const response = await authenticatedFetch(baseUrl + `api/hooks/${selectedFile}`, {
                method: "DELETE",
            });

            if (response.ok) {
                addSuccessToast("删除成功");
                selectedFile = null;
                editor.setValue("");
                await loadFiles();
            } else {
                addErrorToast("删除失败：" + (await response.json()).message);
            }
        });
    }

    async function restartBackend() {
        confirm(`确认要重启服务吗？（只有win平台需要手工重启）`, async () => {
            const response = await authenticatedFetch(baseUrl + `api/hooks/restart`, {
                method: "PUT",
            });

            if (response.ok) {
                addSuccessToast("正在重启……");

                setTimeout(() => {
                    window.location.reload();
                }, 1000);
            } else {
                addErrorToast("重启失败：" + (await response.json()).message);
            }
        });
    }
</script>

<aside class="page-sidebar settings-sidebar">
    <div class="sidebar-content">
        <div class="sidebar-title">云函数</div>

        {#each files as file}
            <button
                class="sidebar-list-item"
                on:click={() => selectFile(file)}
                class:file-selected={file === selectedFile}
            >
                <i class="ri-file-code-line" aria-hidden="true" />
                <span class="txt">{file}</span>
            </button>
        {/each}
    </div>
</aside>

<PageWrapper class="flex-content" showFooter={false}>
    <header class="page-header">
        <div class="btns-group">
            <button type="button" class="btn" on:click={saveFile} disabled={!selectedFile}>
                <i class="ri-save-line" />
                <span class="txt">保存</span>
            </button>
            <button type="button" class="btn btn-outline" on:click={renameFile} disabled={!selectedFile}>
                <i class="ri-braces-line" />
                <span class="txt">重命名</span>
            </button>
            <button type="button" class="btn btn-outline" on:click={deleteFile} disabled={!selectedFile}>
                <i class="ri-delete-bin-5-line" />
                <span class="txt">删除</span>
            </button>
        </div>

        <div class="flex-fill" />

        <div class="btns-group">
            <button type="button" class="btn btn-outline" on:click={restartBackend}>
                <i class="ri-exchange-funds-line" />
                <span class="txt">重启服务</span>
            </button>
            <button type="button" class="btn btn-expanded-sm" on:click={createNewFile}>
                <i class="ri-add-line" />
                <span class="txt">新文件</span>
            </button>
        </div>
    </header>

    <div bind:this={editorContainer} class="editor"></div>

    <div style="">
        <span class="label label-warning">格式化代码：SHIFT + ALT(OPTION) + F</span>
        <span class="label label-warning">Windows版本保存代码之后需要手工重启服务才能生效</span>
    </div>
</PageWrapper>

<style>
    .file-selected {
        background-color: #007bff;
        color: white;
    }

    .file-selected:hover {
        background-color: #007bff;
        color: white;
    }

    .editor {
        flex-grow: 1;
    }
</style>
