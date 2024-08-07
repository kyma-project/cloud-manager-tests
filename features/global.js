
function kymaHasSpecModule(cm, moduleName) {
    if (cm.spec.modules) {
        if (cm.spec.modules.find(function (m) {
            return m.name === moduleName;
        }) !== undefined) {
            return true;
        }
    }
    return false
}


module.exports = {
    steps: [
        {
            rx: `in the kyma (.*) the module (.*) is removed`,
            f: function(ref, moduleName) {
                let cm = load(ref)
                if (kymaHasSpecModule(cm, moduleName)) {
                    cm.spec.modules = cm.spec.modules.filter(function(m) {
                        return m.name !== moduleName;
                    })
                    apply(ref, cm)
                }
            },
        },
        {
            rx: `the kyma (.*) has the module (.*) status (.*)`,
            f: function(ref, moduleName, expectedStatus) {

            }
        }
    ]
}