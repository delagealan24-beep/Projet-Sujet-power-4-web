$base = "http://localhost:8080"
$paths = @("/", "/templates/init", "/templates/play", "/templates/end", "/templates/scoreboard")

Write-Host "Serveur démarré sur $base"
Write-Host "Debug: suivi des pages définies dans ton code..."

Start-Sleep -Seconds 3   

while ($true) {
   
    try {
        $respRoot = Invoke-WebRequest -Uri ($base + "/") -UseBasicParsing -TimeoutSec 5
        $rootContent = $respRoot.Content
    } catch {
        $rootContent = ""
    }

    foreach ($p in $paths) {
        try {
            $resp = Invoke-WebRequest -Uri ($base + $p) -UseBasicParsing -TimeoutSec 5
            $html = $resp.Content

            if ($resp.StatusCode -ne 200) {
                Write-Host ("[{0}] Page:{1} ERREUR HTTP {2}" -f (Get-Date -Format "HH:mm:ss"), $p, $resp.StatusCode)
            }
            elseif ($p -ne "/" -and $html -eq $rootContent) {
                Write-Host ("[{0}] Page:{1} redirige vers / (Status 200)" -f (Get-Date -Format "HH:mm:ss"), $p)
            }
            elseif ([string]::IsNullOrWhiteSpace($html) -or $html.Length -lt 50) {
                Write-Host ("[{0}] Page:{1} Status 200 mais contenu vide/suspect" -f (Get-Date -Format "HH:mm:ss"), $p)
            }
            else {
                Write-Host ("[{0}] Page:{1} OK (Status {2})" -f (Get-Date -Format "HH:mm:ss"), $p, $resp.StatusCode)
            }
        } catch {
            Write-Host ("[{0}] Page:{1} Erreur accès : {2}" -f (Get-Date -Format "HH:mm:ss"), $p, $_.Exception.Message)
        }
    }
    Start-Sleep -Seconds 5
}