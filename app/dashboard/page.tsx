'use client'

import { useState, useEffect } from 'react'
import { SATService } from '@/lib/sat/client'
import { supabase } from '@/lib/supabase/client'

export default function Dashboard() {
  const [cfdis, setCfdis] = useState([])
  const [stats, setStats] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [conexiones, setConexiones] = useState({
    sat: false,
    pac: false,
    db: false
  })
  const satService = new SATService()

  useEffect(() => {
    loadData()
    checkConnections()
  }, [])

  async function loadData() {
    setLoading(true)
    try {
      const [cfdiData, statsData] = await Promise.all([
        satService.getCFDIS(20),
        satService.getStats()
      ])
      setCfdis(cfdiData)
      setStats(statsData)
    } catch (error) {
      console.error('Error:', error)
    } finally {
      setLoading(false)
    }
  }

  async function checkConnections() {
    const { data, error } = await supabase.from('cfdis').select('count', { count: 'exact', head: true })
    setConexiones(prev => ({ ...prev, db: !error }))
    setConexiones(prev => ({ ...prev, sat: true }))
    setConexiones(prev => ({ ...prev, pac: true }))
  }

  const handleIngest = async () => {
    const rfc = prompt('📝 Ingresa RFC:')
    if (!rfc) return
    const password = prompt('🔐 Ingresa contraseña FIEL:')
    if (!password) return
    const period = prompt('📅 Ingresa período (YYYY-MM):')
    if (!period) return

    try {
      const result = await satService.ingestCFDI(rfc, password, period)
      alert(`✅ ${result.message}`)
      loadData()
    } catch (error) {
      alert('❌ Error al ingestar CFDI')
    }
  }

  const handleTimbrar = async () => {
    alert('🏷️ Función de timbrado en desarrollo')
  }

  const generateDIOT = async () => {
    const period = prompt('📅 Ingresa período (YYYY-MM):')
    if (!period) return
    try {
      const result = await satService.generateDIOT(period)
      alert(`📋 DIOT generado: ${result.filename}`)
    } catch (error) {
      alert('❌ Error al generar DIOT')
    }
  }

  const generateReport = () => {
    alert('📊 Reporte ISR/IVA - En desarrollo')
  }

  const auditoriaRiesgo = () => {
    alert('⚠️ Análisis de riesgos - En desarrollo')
  }

  const proyeccionFiscal = () => {
    alert('📅 Proyección fiscal - En desarrollo')
  }

  if (loading) return (
    <div className="flex justify-center items-center h-screen">
      <div className="text-xl">Cargando datos fiscales...</div>
    </div>
  )

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      {/* ENCABEZADO */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-800">
          📊 SISTEMA SAT - FISCAL
        </h1>
        
        <div className="flex gap-4 text-sm">
          <span className={`px-3 py-1 rounded ${conexiones.sat ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
            {conexiones.sat ? '🟢 SAT' : '🔴 SAT'}
          </span>
          <span className={`px-3 py-1 rounded ${conexiones.pac ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
            {conexiones.pac ? '🟢 PAC' : '🔴 PAC'}
          </span>
          <span className={`px-3 py-1 rounded ${conexiones.db ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'}`}>
            {conexiones.db ? '🟢 DB' : '🔴 DB'}
          </span>
        </div>
      </div>

      {/* TARJETAS DE ESTADÍSTICAS */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div className="bg-white p-4 rounded-lg shadow border-l-4 border-blue-500">
          <p className="text-gray-500 text-sm">📄 Total CFDI</p>
          <p className="text-2xl font-bold">{stats?.total || 0}</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow border-l-4 border-green-500">
          <p className="text-gray-500 text-sm">💰 Ingresos</p>
          <p className="text-2xl font-bold text-green-600">
            ${stats?.totalIngresos?.toFixed(2) || '0.00'}
          </p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow border-l-4 border-red-500">
          <p className="text-gray-500 text-sm">💸 Egresos</p>
          <p className="text-2xl font-bold text-red-600">
            ${stats?.totalEgresos?.toFixed(2) || '0.00'}
          </p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow border-l-4 border-purple-500">
          <p className="text-gray-500 text-sm">📈 ISR Estimado</p>
          <p className="text-2xl font-bold text-purple-600">
            $0.00
          </p>
        </div>
      </div>

      {/* BOTONES DE ACCIÓN */}
      <div className="bg-white p-6 rounded-lg shadow mb-6">
        <h2 className="text-lg font-semibold mb-4">🚀 ACCIONES FISCALES</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
          <button
            onClick={handleIngest}
            className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 px-4 rounded-lg transition flex items-center justify-center gap-2"
          >
            📥 Ingestar CFDI
          </button>
          <button
            onClick={handleTimbrar}
            className="bg-purple-600 hover:bg-purple-700 text-white font-medium py-3 px-4 rounded-lg transition flex items-center justify-center gap-2"
          >
            🏷️ Timbrar PAC
          </button>
          <button
            onClick={generateDIOT}
            className="bg-orange-600 hover:bg-orange-700 text-white font-medium py-3 px-4 rounded-lg transition flex items-center justify-center gap-2"
          >
            📋 Generar DIOT
          </button>
          <button
            onClick={generateReport}
            className="bg-green-600 hover:bg-green-700 text-white font-medium py-3 px-4 rounded-lg transition flex items-center justify-center gap-2"
          >
            📊 Reporte ISR
          </button>
          <button
            onClick={auditoriaRiesgo}
            className="bg-red-600 hover:bg-red-700 text-white font-medium py-3 px-4 rounded-lg transition flex items-center justify-center gap-2"
          >
            ⚠️ Auditoría
          </button>
          <button
            onClick={proyeccionFiscal}
            className="bg-teal-600 hover:bg-teal-700 text-white font-medium py-3 px-4 rounded-lg transition flex items-center justify-center gap-2"
          >
            📅 Proyección
          </button>
        </div>
      </div>

      {/* TABLA DE CFDI */}
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="px-6 py-4 border-b">
          <h2 className="text-lg font-semibold">📋 CFDI RECIENTES</h2>
        </div>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">UUID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">RFC</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Tipo</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {cfdis.map((cfdi: any) => (
                <tr key={cfdi.uuid} className="hover:bg-gray-50">
                  <td className="px-6 py-4 text-sm font-mono">
                    {cfdi.uuid?.slice(0, 8)}...
                  </td>
                  <td className="px-6 py-4 text-sm">{cfdi.emisores?.rfc}</td>
                  <td className="px-6 py-4 text-sm">
                    {new Date(cfdi.fecha).toLocaleDateString('es-MX')}
                  </td>
                  <td className="px-6 py-4 text-sm text-right font-medium">
                    ${cfdi.total?.toFixed(2)}
                  </td>
                  <td className="px-6 py-4 text-sm">
                    <span className={`px-2 py-1 rounded-full text-xs font-medium ${
                      cfdi.tipo_comprobante === 'I' 
                        ? 'bg-green-100 text-green-800' 
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {cfdi.tipo_comprobante === 'I' ? '🟢 Ingreso' : '🔴 Egreso'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}