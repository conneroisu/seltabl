import { LanguageClient, LanguageClientOptions, ServerOptions, TransportKind } from 'vscode-languageclient/node';

let client: LanguageClient;

function activate(context) {
    const serverOptions: ServerOptions = {
        command: "seltabls lsp",
        transport: TransportKind.stdio,
        documentSelector: [{ scheme: 'file', language: 'go' }],
    };
    let clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'go' }],
    };
    client = new LanguageClient(
        'seltabls',
        'seltabls',
        serverOptions,
        clientOptions
    );
    client.start();
}

function deactivate() {
    if (!client) {
        return undefined;
    }
    return client.stop();
}

export {
    activate,
    deactivate
};
