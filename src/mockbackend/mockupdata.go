package mockbackend

var serviceV1RootText = []byte(` {
    "@odata.type": "#ServiceRoot.v1_0_2.ServiceRoot",
    "Id": "RootService",
    "Name": "Root Service",
    "RedfishVersion": "1.0.2",
    "UUID": "92384634-2938-2342-8820-489239905423",
    "Systems": {
        "@odata.id": "/redfish/v1/Systems"
    },
    "Chassis": {
        "@odata.id": "/redfish/v1/Chassis"
    },
    "Managers": {
        "@odata.id": "/redfish/v1/Managers"
    },
    "Tasks": {
        "@odata.id": "/redfish/v1/TaskService"
    },
    "SessionService": {
        "@odata.id": "/redfish/v1/SessionService"
    },
    "AccountService": {
        "@odata.id": "/redfish/v1/AccountService"
    },
    "EventService": {
        "@odata.id": "/redfish/v1/EventService"
    },
    "Links": {
        "Sessions": {
            "@odata.id": "/redfish/v1/SessionService/Sessions"
        }
    },
    "Oem": {},
    "@odata.context": "/redfish/v1/$metadata#ServiceRoot",
    "@odata.id": "/redfish/v1/",
    "@Redfish.Copyright": "Copyright 2014-2016 Distributed Management Task Force, Inc. (DMTF). For the full DMTF copyright policy, see http://www.dmtf.org/about/policies/copyright."
}`)

var SystemCollectionText = []byte(`{
    "@odata.type": "#ComputerSystemCollection.ComputerSystemCollection",
    "Name": "Computer System Collection",
    "Members@odata.count": 1,
    "Members": [
        {
            "@odata.type": "#ComputerSystem.v1_1_0.ComputerSystem",
            "Id": "437XR1138R2",
            "Name": "WebFrontEnd483",
            "SystemType": "Physical",
            "AssetTag": "Chicago-45Z-2381",
            "Manufacturer": "Contoso",
            "Model": "3500RX",
            "SKU": "8675309",
            "SerialNumber": "437XR1138R2",
            "PartNumber": "224071-J23",
            "Description": "Web Front End node",
            "UUID": "38947555-7742-3448-3784-823347823834",
            "HostName": "web483",
            "Status": {
                "State": "Enabled",
                "Health": "OK",
                "HealthRollUp": "OK"
            },
            "IndicatorLED": "Off",
            "PowerState": "On",
            "Boot": {
                "BootSourceOverrideEnabled": "Once",
                "BootSourceOverrideTarget": "Pxe",
                "BootSourceOverrideTarget@Redfish.AllowableValues": [
                    "None",
                    "Pxe",
                    "Cd",
                    "Usb",
                    "Hdd",
                    "BiosSetup",
                    "Utilities",
                    "Diags",
                    "SDCard",
                    "UefiTarget"
                ],
                "BootSourceOverrideMode": "UEFI",
                "UefiTargetBootSourceOverride": "/0x31/0x33/0x01/0x01"
            },
            "TrustedModules": [
                {
                    "FirmwareVersion": "1.13b",
                    "InterfaceType": "TPM1_2",
                    "Status": {
                        "State": "Enabled",
                        "Health": "OK"
                    }
                }
            ],
            "Oem": {
                "Contoso": {
                    "@odata.type": "http://Contoso.com/Schema#Contoso.ComputerSystem",
                    "ProductionLocation": {
                        "FacilityName": "PacWest Production Facility",
                        "Country": "USA"
                    }
                },
                "Chipwise": {
                    "@odata.type": "http://Chipwise.com/Schema#Chipwise.ComputerSystem",
                    "Style": "Executive"
                }
            },
            "BiosVersion": "P79 v1.33 (02/28/2015)",
            "ProcessorSummary": {
                "Count": 2,
                "ProcessorFamily": "Multi-Core Intel(R) Xeon(R) processor 7xxx Series",
                "Status": {
                    "State": "Enabled",
                    "Health": "OK",
                    "HealthRollUp": "OK"
                }
            },
            "MemorySummary": {
                "TotalSystemMemoryGiB": 96,
                "Status": {
                    "State": "Enabled",
                    "Health": "OK",
                    "HealthRollUp": "OK"
                }
            },
            "Bios": {
                "@odata.id": "/redfish/v1/Systems/437XR1138R2/BIOS"
            },
            "Processors": {
                "@odata.id": "/redfish/v1/Systems/437XR1138R2/Processors"
            },
            "Memory": {
                "@odata.id": "/redfish/v1/Systems/437XR1138R2/Memory"
            },
            "EthernetInterfaces": {
                "@odata.id": "/redfish/v1/Systems/437XR1138R2/EthernetInterfaces"
            },
            "SimpleStorage": {
                "@odata.id": "/redfish/v1/Systems/437XR1138R2/SimpleStorage"
            },
            "LogServices": {
                "@odata.id": "/redfish/v1/Systems/437XR1138R2/LogServices"
            },
            "Links": {
                "Chassis": [
                    {
                        "@odata.id": "/redfish/v1/Chassis/1U"
                    }
                ],
                "ManagedBy": [
                    {
                        "@odata.id": "/redfish/v1/Managers/BMC"
                    }
                ]
            },
            "Actions": {
                "#ComputerSystem.Reset": {
                    "target": "/redfish/v1/Systems/437XR1138R2/Actions/ComputerSystem.Reset",
                    "ResetType@Redfish.AllowableValues": [
                        "On",
                        "ForceOff",
                        "GracefulShutdown",
                        "GracefulRestart",
                        "ForceRestart",
                        "Nmi",
                        "ForceOn",
                        "PushPowerButton"
                    ]
                },
                "Oem": {
                    "#Contoso.Reset": {
                        "target": "/redfish/v1/Systems/437XR1138R2/Oem/Contoso/Actions/Contoso.Reset"
                    }
                }
            },
            "@odata.context": "/redfish/v1/$metadata#ComputerSystem.ComputerSystem",
            "@odata.id": "/redfish/v1/Systems/437XR1138R2",
            "@Redfish.Copyright": "Copyright 2014-2016 Distributed Management Task Force, Inc. (DMTF). For the full DMTF copyright policy, see http://www.dmtf.org/about/policies/copyright."
        }
    ],
    "@odata.context": "/redfish/v1/$metadata#Systems",
    "@odata.id": "/redfish/v1/Systems",
    "@Redfish.Copyright": "Copyright 2014-2016 Distributed Management Task Force, Inc. (DMTF). For the full DMTF copyright policy, see http://www.dmtf.org/about/policies/copyright."
}`)