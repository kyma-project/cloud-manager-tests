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


  @aws @allShoots @dev
  Scenario: Installed CRDs
    When CRDs are loaded
    Then CRDs exist:
      | IpRange          |
      | AwsNfsVolume     |
      | AwsVpcPeering    |
      | AwsRedisInstance |
    And CRDs do not exist:
      | GcpNfsVolume        |
      | GcpNfsVolumeBackup  |
      | GcpNfsVolumeRestore |
      | GcpNfsBackupSchedule|
      | GcpRedisInstance    |
      | GcpVpcPeering       |
      | AzureVpcPeering     |
      | AzureRedisInstance  |

  @gcp @allShoots @dev
  Scenario: Installed CRDs
    When CRDs are loaded
    Then CRDs exist:
      | IpRange             |
      | GcpNfsVolume        |
      | GcpNfsVolumeBackup  |
      | GcpNfsVolumeRestore |
      | GcpNfsBackupSchedule|
      | GcpRedisInstance    |
      | GcpVpcPeering       |
    And CRDs do not exist:
      | AwsNfsVolume       |
      | AwsVpcPeering      |
      | AwsRedisInstance   |
      | AzureVpcPeering    |
      | AzureRedisInstance |

   @azure @allShoots @dev
     Scenario: Installed CRDs
       When CRDs are loaded
       Then CRDs exist:
         | IpRange            |
         | AzureVpcPeering    |
         | AzureRedisInstance |
       And CRDs do not exist:
         | AwsNfsVolume       |
         | AwsVpcPeering      |
         | AwsRedisInstance   |
         | GcpNfsVolume        |
         | GcpNfsVolumeBackup  |
         | GcpNfsVolumeRestore |
         | GcpNfsBackupSchedule|
         | GcpRedisInstance    |
         | GcpVpcPeering       |
