const vscode = require('vscode');
const path = require('path');
const fs = require('fs');

class DuckDocumentLinkProvider {
    provideDocumentLinks(document, token) {
        const links = [];
        const text = document.getText();
        const regex = /import\s+.*?from\s+["']([^"']+)["']/g;
        let match;

        while ((match = regex.exec(text)) !== null) {
            const importPath = match[1];
            const startOffset = match.index + match[0].indexOf(importPath);
            const endOffset = startOffset + importPath.length;
            const startPos = document.positionAt(startOffset);
            const endPos = document.positionAt(endOffset);
            const range = new vscode.Range(startPos, endPos);

            let targetPath = importPath;
            if (!path.extname(targetPath)) {
                targetPath += '.duck';
            }

            const currentDir = path.dirname(document.uri.fsPath);
            const resolvedPath = path.resolve(currentDir, targetPath);
            const targetUri = vscode.Uri.file(resolvedPath);

            links.push(new vscode.DocumentLink(range, targetUri));
        }
        return links;
    }
}

class DuckDefinitionProvider {
    provideDefinition(document, position, token) {
        const wordRange = document.getWordRangeAtPosition(position);
        if (!wordRange) {
            return null;
        }
        const word = document.getText(wordRange);
        const text = document.getText();
        const currentDir = path.dirname(document.uri.fsPath);

        // 1. Check if defined in current file
        const localDeclRegex = new RegExp(`\\b(class|service|func)\\s+${word}\\b`);
        const localMatch = text.match(localDeclRegex);
        if (localMatch) {
            const index = text.indexOf(localMatch[0]);
            const pos = document.positionAt(index + localMatch[0].indexOf(word));
            return new vscode.Location(document.uri, pos);
        }

        // 2. Check if the word is imported from another file
        const importRegex = new RegExp(`import\\s+{[^}]*\\b${word}\\b[^}]*}\\s+from\\s+["']([^"']+)["']`);
        const importMatch = text.match(importRegex);
        if (importMatch) {
            const importPath = importMatch[1];
            const targetUri = resolveImportPath(currentDir, importPath);
            if (targetUri && fs.existsSync(targetUri.fsPath)) {
                const targetText = fs.readFileSync(targetUri.fsPath, 'utf8');
                const targetDeclRegex = new RegExp(`\\b(class|service|func)\\s+${word}\\b`);
                const targetMatch = targetText.match(targetDeclRegex);
                if (targetMatch) {
                    const index = targetText.indexOf(targetMatch[0]);
                    const offset = index + targetMatch[0].indexOf(word);
                    const pos = getPositionFromOffset(targetText, offset);
                    return new vscode.Location(targetUri, pos);
                }
            }
        }

        // 3. Check for method definition, e.g. service.getHealth() or this.service.getHealth()
        const line = document.lineAt(position.line).text;
        const methodCallRegex = new RegExp(`(?:this\\.)?([a-zA-Z0-9_]+)\\.${word}\\b`);
        const methodMatch = line.match(methodCallRegex);
        if (methodMatch) {
            const varName = methodMatch[1];
            
            const varDeclRegex = new RegExp(`\\b${varName}\\s*:\\s*([a-zA-Z0-9_]+)`);
            const varDeclMatch = text.match(varDeclRegex);
            if (varDeclMatch) {
                const typeName = varDeclMatch[1];
                
                const typeImportRegex = new RegExp(`import\\s+{[^}]*\\b${typeName}\\b[^}]*}\\s+from\\s+["']([^"']+)["']`);
                const typeImportMatch = text.match(typeImportRegex);
                if (typeImportMatch) {
                    const importPath = typeImportMatch[1];
                    const targetUri = resolveImportPath(currentDir, importPath);
                    if (targetUri && fs.existsSync(targetUri.fsPath)) {
                        const targetText = fs.readFileSync(targetUri.fsPath, 'utf8');
                        const methodRegex = new RegExp(`\\b${word}\\s*\\(`);
                        const methodMatchInFile = targetText.match(methodRegex);
                        if (methodMatchInFile) {
                            const index = targetText.indexOf(methodMatchInFile[0]);
                            const pos = getPositionFromOffset(targetText, index);
                            return new vscode.Location(targetUri, pos);
                        }
                    }
                }
            }
        }

        return null;
    }
}

function resolveImportPath(currentDir, importPath) {
    let targetPath = importPath;
    if (!path.extname(targetPath)) {
        targetPath += '.duck';
    }
    const resolvedPath = path.resolve(currentDir, targetPath);
    return vscode.Uri.file(resolvedPath);
}

function getPositionFromOffset(text, offset) {
    const lines = text.slice(0, offset).split('\n');
    const line = lines.length - 1;
    const character = lines[lines.length - 1].length;
    return new vscode.Position(line, character);
}

function activate(context) {
    const selector = { language: 'goduck', scheme: 'file' };

    context.subscriptions.push(
        vscode.languages.registerDocumentLinkProvider(selector, new DuckDocumentLinkProvider())
    );

    context.subscriptions.push(
        vscode.languages.registerDefinitionProvider(selector, new DuckDefinitionProvider())
    );
}

function deactivate() {}

module.exports = {
    activate,
    deactivate
};
