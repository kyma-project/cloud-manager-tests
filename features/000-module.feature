Feature: Module enable feature

  @all
  Scenario: Module enable scenario
    Given resource declaration:
      | kyma | Kyma           | "default" | "kyma-system" |
      | cm   | CloudResources | "default" | "kyma-system" |
    Given there are no cloud resources
    And module is removed
    When module is added
    Then eventually value load("cm").status.state equals "Ready"


  @aws
  Scenario: Installed CRDs
    When CRDs are loaded
    Then CRDs exist:
      | IpRange          |               |
      | AwsNfsVolume     |               |
      | AwsVpcPeering    | env == "dev"  |
      | AwsRedisInstance | env == "dev"  |

  @gcp
  Scenario: Installed CRDs
    When CRDs are loaded
    Then CRDs exist:
      | IpRange             |               |
      | GcpNfsVolume        |               |
      | GcpNfsVolumeBackup  | env == "dev"  |
      | GcpNfsVolumeRestore | env == "dev"  |
      | GcpRedisInstance    | env == "dev"  |
      | GcpVpcPeering       | env == "dev"  |
