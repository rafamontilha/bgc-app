Write-Host "API /healthz";      iwr http://api.bgc.local/healthz
Write-Host "API TAM 84 2024-25"; iwr "http://api.bgc.local/market/size?metric=TAM&year_from=2024&year_to=2025&ncm_chapter=84"
Write-Host "API SOM aggressive";  iwr "http://api.bgc.local/market/size?metric=SOM&scenario=aggressive&year_from=2024&year_to=2025"
Write-Host "API routes compare";  iwr "http://api.bgc.local/routes/compare?from=USA&alts=CHN,ARE,IND&ncm_chapter=84&year=2024&tariff_scenario=tarifa10"
