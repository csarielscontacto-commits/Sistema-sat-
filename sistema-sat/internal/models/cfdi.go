package models

import "encoding/xml"

type CFDI struct {
    XMLName          xml.Name    `xml:"cfdi:Comprobante" json:"-"`
    Version          string      `xml:"Version,attr" json:"version"`
    Serie            string      `xml:"Serie,attr" json:"serie"`
    Folio            string      `xml:"Folio,attr" json:"folio"`
    Fecha            string      `xml:"Fecha,attr" json:"fecha"`
    FormaPago        string      `xml:"FormaPago,attr" json:"formaPago"`
    SubTotal         float64     `xml:"SubTotal,attr" json:"subTotal"`
    Descuento        float64     `xml:"Descuento,attr" json:"descuento"`
    Moneda           string      `xml:"Moneda,attr" json:"moneda"`
    Total            float64     `xml:"Total,attr" json:"total"`
    TipoDeComprobante string    `xml:"TipoDeComprobante,attr" json:"tipoDeComprobante"`
    MetodoPago       string      `xml:"MetodoPago,attr" json:"metodoPago"`
    LugarExpedicion  string      `xml:"LugarExpedicion,attr" json:"lugarExpedicion"`
    Emisor           Emisor      `xml:"cfdi:Emisor" json:"emisor"`
    Receptor         Receptor    `xml:"cfdi:Receptor" json:"receptor"`
    Conceptos        []Concepto  `xml:"cfdi:Conceptos>cfdi:Concepto" json:"conceptos"`
    Impuestos        Impuestos   `xml:"cfdi:Impuestos" json:"impuestos"`
    Complemento      Complemento `xml:"cfdi:Complemento" json:"complemento"`
    UUID             string      `json:"uuid"`
    RFC              string      `json:"rfc"`
}

type Emisor struct {
    RFC           string `xml:"Rfc,attr" json:"rfc"`
    Nombre        string `xml:"Nombre,attr" json:"nombre"`
    RegimenFiscal string `xml:"RegimenFiscal,attr" json:"regimenFiscal"`
}

type Receptor struct {
    RFC     string `xml:"Rfc,attr" json:"rfc"`
    Nombre  string `xml:"Nombre,attr" json:"nombre"`
    UsoCFDI string `xml:"UsoCFDI,attr" json:"usoCFDI"`
}

type Concepto struct {
    ClaveProdServ string  `xml:"ClaveProdServ,attr" json:"claveProdServ"`
    Cantidad      float64 `xml:"Cantidad,attr" json:"cantidad"`
    ClaveUnidad   string  `xml:"ClaveUnidad,attr" json:"claveUnidad"`
    Descripcion   string  `xml:"Descripcion,attr" json:"descripcion"`
    ValorUnitario float64 `xml:"ValorUnitario,attr" json:"valorUnitario"`
    Importe       float64 `xml:"Importe,attr" json:"importe"`
}

type Impuestos struct {
    TotalImpuestosTrasladados float64    `xml:"TotalImpuestosTrasladados,attr" json:"totalImpuestosTrasladados"`
    Traslados                 []Traslado `xml:"cfdi:Traslados>cfdi:Traslado" json:"traslados"`
}

type Traslado struct {
    Base       float64 `xml:"Base,attr" json:"base"`
    Impuesto   string  `xml:"Impuesto,attr" json:"impuesto"`
    TipoFactor string  `xml:"TipoFactor,attr" json:"tipoFactor"`
    TasaOCuota float64 `xml:"TasaOCuota,attr" json:"tasaOCuota"`
    Importe    float64 `xml:"Importe,attr" json:"importe"`
}

type Complemento struct {
    TimbreFiscal TimbreFiscal `xml:"tfd:TimbreFiscalDigital" json:"timbreFiscal"`
}

type TimbreFiscal struct {
    UUID             string `xml:"UUID,attr" json:"uuid"`
    Fecha            string `xml:"FechaTimbrado,attr" json:"fechaTimbrado"`
    SelloCFD         string `xml:"SelloCFD,attr" json:"selloCFD"`
    NoCertificadoSAT string `xml:"NoCertificadoSAT,attr" json:"noCertificadoSAT"`
}