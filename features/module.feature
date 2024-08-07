Feature: Module enable feature

  Scenario: Module enable scenario
    Given resource declaration:
      | cm | CloudResources | default | kyma-system |
    Given in the kyma cm the module cloud-manager is removed
    And the kyma cm has the module cloud-manager status None