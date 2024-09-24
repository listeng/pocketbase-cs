<script>
    import { onMount } from "svelte";
    import * as monaco from "monaco-editor";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";
    import { confirm } from "@/stores/confirmation";
    import PageWrapper from "@/components/base/PageWrapper.svelte";

    import editorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker";
    import tsWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker";

    let files = [];
    let selectedFile = null;
    let editor;
    let editorContainer;
    const baseUrl = import.meta.env.PB_BACKEND_URL;

    async function loadDtsFile(url) {
    const response = await fetch(url);
    const dtsContent = await response.text();
    monaco.languages.typescript.javascriptDefaults.addExtraLib(dtsContent, url);
  }

    onMount(async () => {
        self.MonacoEnvironment = {
            getWorker: function (_, label) {
                if (label === "typescript" || label === "javascript") {
                    return new tsWorker();
                }
                return new editorWorker();
            },
        };

        loadFiles();

        editor = monaco.editor.create(editorContainer, {
            value: "",
            language: "javascript",
            theme: "vs-dark",
        });

        await loadDtsFile(baseUrl + 'types.d.ts');

        return () => {
            editor.dispose();
        };
    });

    async function authenticatedFetch(url, options = {}) {
        const authData = localStorage.getItem("pb_admin_auth");
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
        files = await response.json();
    }

    async function selectFile(filename) {
        selectedFile = filename;
        const response = await authenticatedFetch(baseUrl + `api/hooks/${filename}`);
        const content = await response.json();
        editor.setValue(content);
    }

    async function saveFile() {
        if (!selectedFile) return;

        const content = editor.getValue();
        await authenticatedFetch(baseUrl + `api/hooks/${selectedFile}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ content }),
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

<PageWrapper class="flex-content">
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
            <button type="button" class="btn btn-expanded-sm" on:click={createNewFile}>
                <i class="ri-add-line" />
                <span class="txt">新文件</span>
            </button>
        </div>
    </header>

    <div bind:this={editorContainer} class="editor"></div>
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
