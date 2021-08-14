; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!

#define MyAppName "OctalSpan"
#define MyAppVersion "0.1.3"
#define MyAppPublisher "Cutrer Technologies"
#define MyAppURL "https://www.cutrertech.com"
#define MyAppExeName "octalspan.syslogd.exe"

[Setup]
; NOTE: The value of AppId uniquely identifies this application. Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={{EB1DC14A-D82A-4BBB-B185-6F2B448B18CC}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
;AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\octalspan
DisableDirPage=yes
DisableProgramGroupPage=yes
LicenseFile=LICENSE
; Uncomment the following line to run in non administrative install mode (install for current user only.)
;PrivilegesRequired=lowest
OutputDir=dist
OutputBaseFilename=octalspan-{#MyAppVersion}
Compression=lzma
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "bin\octalspan.syslogd.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\octalspan.yml"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\octalspan.svc.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\octalspan.svc.xml"; DestDir: "{app}"; Flags: ignoreversion
Source: "bin\LICENSE"; DestDir: "{app}"; Flags: ignoreversion
; Source: "C:\Users\joncu\OneDrive\dev\golang-apps\octalspan\bin\*"; DestDir: "{app}"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

