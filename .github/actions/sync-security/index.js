const core = require('@actions/core');
const github = require('@actions/github');
const fsAsync = require('node:fs/promises');
const process = require('node:process');
const util = require('node:util');
const path = require('node:path');
const child_process = require('node:child_process');
const exec = util.promisify(child_process.exec);

function getInputOwner() {
    let result = core.getInput('owner');
    if (result) {
        return result;
    }
    const {owner} = github.context.repo;
    if (!owner) {
        throw new Error('Unable to determine owner - not specified as input nor able to read it from context');
    }
    return owner;
}

function getInputRepo() {
    let result = core.getInput('repo');
    if (result) {
        return result;
    }
    const {repo} = github.context.repo;
    if (!repo) {
        throw new Error('Unable to determine repo - not specified as input nor able to read it from context');
    }
    return repo;
}

async function getRootDir() {
    let result = process.env.GITHUB_WORKSPACE;
    if (!result) {
        result = await fsAsync.realpath(path.join(__dirname, '..', '..', '..'));
    }
    return result;
}

function toStr(strOrBuffer) {
    if (!strOrBuffer) {
        return '';
    }
    if (typeof strOrBuffer === 'string') {
        return strOrBuffer;
    }
    if (typeof strOrBuffer === 'object' && typeof strOrBuffer.toString === 'function') {
        return strOrBuffer.toString();
    }
    return '';
}

let dryRun = false;

async function run() {
    dryRun = core.getBooleanInput('dryRun');
    const rootDir = await getRootDir();
    const owner = getInputOwner();
    const repo = getInputRepo();

    const projectNum = core.getInput('project_num');
    if (!projectNum) {
        throw new Error('Input project_num is required');
    }

    const externalIssueLinkFieldName = core.getInput('externalIssueLinkFieldName') || 'External Issue Link';

    const issuesFilename = core.getInput('issues');
    if (!issuesFilename) {
        throw new Error('Input issues is required');
    }
    const issues = JSON.parse(await fsAsync.readFile(path.join(rootDir, issuesFilename)));

    const project = await loadProject(rootDir);
    const items = project.items.nodes;

    console.log(`Issues loaded: ${issues.length}`)
    console.log(util.inspect(issues, undefined, 5));
    console.log('--------------------------------------------');
    console.log(`Items loaded: ${items.length}`);
    console.log(util.inspect(items, undefined, 5));
    console.log('--------------------------------------------');

    const fields = await listFields(projectNum, owner);
    const externalIssueLinkField = findFieldByName(fields, externalIssueLinkFieldName);

    for (let issue of issues) {
        console.log(`Issue ${issue.title}`);
        const item = findItemByTitle(items, issue.title);
        if (!item) {
            await createItem(projectNum, owner, issue);
            continue;
        }

        if (getItemFieldValue(item, externalIssueLinkFieldName) !== issue.url) {
            await updateItemFieldText(project.id, item.id, externalIssueLinkField.id, issue.url);
        }
    }
}

async function loadProject(rootDir) {
    const itemsFilename = core.getInput('items');
    if (!itemsFilename) {
        throw new Error('Input items is required');
    }

    const content = JSON.parse(await fsAsync.readFile(path.join(rootDir, itemsFilename)));

    let project = undefined;
    for (let chunk of content) {
        if (!project) {
            project = chunk.data.repository.projectV2;
        } else {
            project.items.nodes.push(...chunk.data.repository.projectV2.items.nodes);
        }
    }
    return project;
}

function findFieldByName(fields, name) {
    for (let field of fields) {
        if (field.name === name) {
            return field;
        }
    }
}

async function listFields(projectNum, owner) {
    const cmd = `gh project field-list ${projectNum} --owner ${owner} --format json`;
    const {stdout, stderr} = await exec(cmd, {
        maxBuffer: 1024 * 1024 * 1024,
    });
    const txt = toStr(stdout);
    return JSON.parse(txt).fields;
}

async function updateItemFieldText(projectId, itemId, fieldId, text) {
    const cmd = `gh project item-edit --project-id ${projectId} --id ${itemId} --field-id ${fieldId} --text "${text.replaceAll('"', '\\""')}"`;
    console.log(`Updating project ${projectId} item ${itemId} field ${fieldId} text ${text}`);
    console.log(cmd);
    if (!dryRun) {
        const {stdout, stderr} = await exec(cmd, {
            maxBuffer: 1024 * 1024 * 1024,
        });
        console.log(toStr(stdout));
    }
}


async function createItem(projectNum, owner, issue) {
    const title = issue.title.replaceAll('"', '\\"');
    const body = `
# ${issue.title}
${issue.url}
`;
    console.log(`Creating item: ${title}`);
    const cmd = `gh project item-create ${projectNum} --owner "${owner}" --title "${title}" --body "${body}"`
    console.log(cmd);
    if (!dryRun) {
        const {stdout, stderr} = await exec(cmd, {
            maxBuffer: 1024 * 1024 * 1024,
        });
        console.log(toStr(stdout));
    }
}

function findItemByTitle(items, title) {
    for (let item of items) {
        if (item.content?.title === title) {
            return item;
        }
    }
    return undefined;
}

function getItemFieldValue(item, fieldName) {
    for (let field of item.fieldValues.nodes) {
        if (field.field.name === fieldName) {
            return field.text;
        }
    }
    return undefined;
}

async function main() {
    try {
        run();
    } catch (err) {
        core.setFailed(err.message);
    }
}

process.env.INPUT_ISSUES = '../../../tmp/issues.json';
process.env.INPUT_ITEMS = '../../../tmp/items.json';
process.env.INPUT_EXTERNALISSUELINKFIELDNAME = 'External Issue Link';
process.env.INPUT_PROJECT_NUM = '55';
process.env.INPUT_OWNER = 'kyma-project';
process.env.INPUT_REPO = 'cloud-manager-tests';

main();