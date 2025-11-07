# Certificados ICP-Brasil

Este diretório armazena certificados digitais ICP-Brasil para autenticação com sistemas governamentais.

## ⚠️ IMPORTANTE - SEGURANÇA

**NUNCA commite certificados reais neste repositório!**

- Este diretório está protegido por `.gitignore`
- Certificados devem ser gerenciados via Certificate Manager Service
- Em produção, use Kubernetes Secrets ou HashiCorp Vault
- Use apenas certificados de sandbox/homologação em desenvolvimento

## Tipos de Certificados Suportados

### A1 - Certificado em Software
- **Formato**: PFX/P12
- **Armazenamento**: Arquivo criptografado
- **Validade**: Típicamente 1 ano
- **Uso**: Desenvolvimento e sistemas de menor criticidade

### A3 - Certificado em Hardware
- **Formato**: Token USB ou Smart Card (HSM)
- **Armazenamento**: Hardware criptográfico
- **Validade**: Típicamente 3-5 anos
- **Uso**: Produção e sistemas críticos

## Configuração Local (Desenvolvimento)

### Certificado A1

1. Copie seu certificado `.pfx` para este diretório
2. Configure as variáveis de ambiente:

```bash
export ICP_CERT_TYPE=A1
export ICP_CERT_PATH=/app/certs/seu-certificado.pfx
export ICP_CERT_PASSWORD=sua-senha-segura
```

### Certificado A3

1. Instale os drivers do fabricante do token/smartcard
2. Configure as variáveis de ambiente:

```bash
export ICP_CERT_TYPE=A3
export ICP_CERT_SLOT=0
export ICP_CERT_PIN=seu-pin
export ICP_PKCS11_LIB=/usr/lib/libaetpkss.so  # Caminho da biblioteca PKCS#11
```

## Gerenciamento via Certificate Manager

O Certificate Manager Service centraliza a gestão de todos os certificados:

- Auto-renovação antes da expiração
- Validação contínua
- Audit trail completo
- Alertas de expiração
- Multi-tenancy

Ver: `services/cert-manager/`

## Referências

- [ICP-Brasil](https://www.gov.br/iti/pt-br/assuntos/icp-brasil)
- [PKCS#11 Specification](https://docs.oasis-open.org/pkcs11/pkcs11-base/v2.40/pkcs11-base-v2.40.html)
- Documentação interna: `docs/CERTIFICATE-MANAGEMENT.md`
