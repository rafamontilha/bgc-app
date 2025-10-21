# Start API with local PostgreSQL
$env:DB_HOST = "localhost"
$env:DB_USER = "postgres"
$env:DB_PASS = "postgres"
$env:DB_NAME = "bgc"

Set-Location "C:\Users\rafae\OneDrive\Documentos\Projetos\Brasil Global Conect\bgc-app\api"
& "C:\Program Files\Go\bin\go.exe" run cmd/api/main.go
