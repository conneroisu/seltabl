// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
const vscode = require('vscode');

// This method is called when your extension is activated
// Your extension is activated the very first time the command is executed

const path = require('path');
const { workspace, window } = require('vscode');
const { LanguageClient, TransportKind } = require('vscode-languageclient/node');

let client;
/**
 * @param {vscode.ExtensionContext} context
 */
function activate(context) {
    let serverModule = context.asAbsolutePath(
        path.join('seltabls', 'lsp')
    );

    let debugOptions = { execArgv: ['--nolazy', '--inspect=6009'] };

    let serverOptions = {
        run: { module: serverModule, transport: TransportKind.stdio },
        debug: {
            module: serverModule,
            transport: TransportKind.stdio,
            options: debugOptions
        }
    };

    let clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'go' }],
        synchronize: {
            fileEvents: workspace.createFileSystemWatcher('**/.clientrc')
        }
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

module.exports = {
    activate,
    deactivate
};
