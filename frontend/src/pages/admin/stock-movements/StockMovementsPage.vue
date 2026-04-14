<template>
  <div class="space-y-6">
<div class="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h2 class="text-xl font-bold text-slate-900">Mutasi Stok</h2>
          <p class="text-sm text-slate-500">Lacak pergerakan stok masuk, keluar, dan penyesuaian.</p>
        </div>
        <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
          <select v-model="productFilter" class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
            <option value="">Semua Produk</option>
            <option v-for="product in products" :key="product.id" :value="product.id">{{ product.name }}</option>
          </select>
          <select v-model="warehouseFilter" class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
            <option value="">Semua Gudang</option>
            <option v-for="warehouse in warehouses" :key="warehouse.id" :value="warehouse.id">{{ warehouse.name }}</option>
          </select>
          <input v-model="movementTypeFilter" type="text" placeholder="Filter tipe mutasi" class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
          <button @click="fetchStockMovements" class="rounded-full border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-50 transition">Segarkan</button>
        </div>
    </div>

    <div v-if="stockMovementError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ stockMovementError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Tipe</th>
            <th class="px-6 py-4">Produk</th>
            <th class="px-6 py-4">Gudang</th>
            <th class="px-6 py-4">Kuantitas</th>
            <th class="px-6 py-4">Referensi</th>
            <th class="px-6 py-4">Catatan</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="movement in stockMovements" :key="movement.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ movement.movementType }}</td>
            <td class="px-6 py-4">{{ getProductName(movement.productId) }}</td>
            <td class="px-6 py-4">{{ getWarehouseName(movement.warehouseId) }}</td>
            <td class="px-6 py-4">{{ movement.quantity }}</td>
            <td class="px-6 py-4">{{ movement.referenceType }} #{{ movement.referenceId }}</td>
            <td class="px-6 py-4">{{ movement.notes }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface StockMovement {
  id: number
  productId: number
  warehouseId: number
  movementType: string
  quantity: number
  referenceType: string
  referenceId: string
  notes: string
}

interface Product {
  id: number
  name: string
}

interface Warehouse {
  id: number
  name: string
}

const router = useRouter()
const stockMovements = ref<StockMovement[]>([])
const products = ref<Product[]>([])
const warehouses = ref<Warehouse[]>([])
const productFilter = ref<number | ''>('')
const warehouseFilter = ref<number | ''>('')
const movementTypeFilter = ref('')
const stockMovementError = ref('')

const fetchStockMovements = async () => {
  stockMovementError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const params = new URLSearchParams()
    if (productFilter.value !== '') params.set('productId', String(productFilter.value))
    if (warehouseFilter.value !== '') params.set('warehouseId', String(warehouseFilter.value))
    if (movementTypeFilter.value.trim()) params.set('movementType', movementTypeFilter.value.trim())
    const url = `/api/stock_movements${params.toString() ? `?${params.toString()}` : ''}`
    const response = await fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    stockMovements.value = result.data || []
  } catch {
    stockMovementError.value = 'Gagal memuat mutasi stok'
  }
}

const fetchProducts = async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/products', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    products.value = result.data || []
  } catch {
    // ignore
  }
}

const fetchWarehouses = async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/warehouses', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    warehouses.value = result.data || []
  } catch {
    // ignore
  }
}

const getProductName = (productId: number) => products.value.find((item) => item.id === productId)?.name || `#${productId}`
const getWarehouseName = (warehouseId: number) => warehouses.value.find((item) => item.id === warehouseId)?.name || `#${warehouseId}`

onMounted(async () => {
  await Promise.all([fetchStockMovements(), fetchProducts(), fetchWarehouses()])
})
</script>

<style scoped>
</style>
