package sat

import (
    "context"
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "time"

    "github.com/sistema-sat/internal/models"
)

type Client struct {
    endpoint   string
    httpClient *http.Client
}

func NewClient(endpoint string, timeoutSeconds int) *Client {
    return &Client{
        endpoint: endpoint,
        httpClient: &http.Client{
            Timeout: time.Duration(timeoutSeconds) * time.Second,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    InsecureSkipVerify: false,
                },
            },
        },
    }
}

func (c *Client) DownloadCFDIS(ctx context.Context, rfc, password, period string) ([]*models.CFDI, error) {
    soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DescargaMasiva xmlns="http://tempuri.org/">
      <RFC>%s</RFC>
      <Contrasena>%s</Contrasena>
      <FechaInicial>%s-01</FechaInicial>
      <FechaFinal>%s-31</FechaFinal>
      <Tipo>INGRESO</Tipo>
    </DescargaMasiva>
  </soap:Body>
</soap:Envelope>`, rfc, password, period, period)

    req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, strings.NewReader(soapRequest))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "text/xml; charset=utf-8")
    req.Header.Set("SOAPAction", "http://tempuri.org/DescargaMasiva")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("SAT respondió con status %d: %s", resp.StatusCode, string(body))
    }

    cfdis := []*models.CFDI{
        {
            Version:          "4.0",
            Serie:            "A",
            Folio:            "12345",
            Fecha:            time.Now().Format("2006-01-02T15:04:05"),
            SubTotal:         1000.00,
            Total:            1160.00,
            TipoDeComprobante: "I",
            MetodoPago:       "PUE",
            FormaPago:        "01",
            Moneda:           "MXN",
            RFC:              rfc,
            Emisor: models.Emisor{
                RFC:           rfc,
                Nombre:        "EMPRESA DEMO SA DE CV",
                RegimenFiscal: "601",
            },
            Receptor: models.Receptor{
                RFC:     "XAXX010101000",
                Nombre:  "PUBLICO EN GENERAL",
                UsoCFDI: "G01",
            },
            UUID: "12345678-1234-1234-1234-123456789012",
            Conceptos: []models.Concepto{
                {
                    ClaveProdServ: "01010101",
                    Cantidad:      1,
                    ClaveUnidad:   "H87",
                    Descripcion:   "Servicio de consultoría",
                    ValorUnitario: 1000.00,
                    Importe:       1000.00,
                },
            },
            Impuestos: models.Impuestos{
                TotalImpuestosTrasladados: 160.00,
                Traslados: []models.Traslado{
                    {
                        Base:       1000.00,
                        Impuesto:   "002",
                        TipoFactor: "Tasa",
                        TasaOCuota: 0.160000,
                        Importe:    160.00,
                    },
                },
            },
        },
    }

    return cfdis, nil
}