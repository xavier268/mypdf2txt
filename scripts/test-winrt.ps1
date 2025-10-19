# Test script pour débugger les API WinRT

Add-Type -AssemblyName System.Runtime.WindowsRuntime
[Windows.Storage.StorageFile,Windows.Storage,ContentType=WindowsRuntime] | Out-Null

$task = [Windows.Storage.StorageFile]::GetFileFromPathAsync("C:\Windows\System32\notepad.exe")

Write-Host "Type complet: $($task.GetType().FullName)"
Write-Host "GenericTypeArguments count: $($task.GetType().GenericTypeArguments.Count)"
Write-Host "GenericTypeArguments: $($task.GetType().GenericTypeArguments)"
Write-Host "GetGenericArguments count: $($task.GetType().GetGenericArguments().Count)"
Write-Host "GetGenericArguments: $($task.GetType().GetGenericArguments())"

# Tester si c'est un IAsyncOperation
$interfaces = $task.GetType().GetInterfaces()
Write-Host "Interfaces:"
foreach ($interface in $interfaces) {
    Write-Host "  - $($interface.FullName)"
    if ($interface.IsGenericType) {
        Write-Host "    Generic args: $($interface.GetGenericArguments())"
    }
}
