<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen Inventaris</h2>
        <p class="text-sm text-slate-500">Pantau dan perbarui stok produk di gudang.</p>
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
        <button @click="fetchInventory" class="rounded-full border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-50 transition">Segarkan</button>
      </div>
    </div>

    <div v-if="inventoryError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ inventoryError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Produk</th>
            <th class="px-6 py-4">Gudang</th>
            <th class="px-6 py-4">Qty On Hand</th>
            <th class="px-6 py-4">Qty Reserved</th>
            <th class="px-6 py-4">Available</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="item in inventory" :key="item.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ getProductName(item.productId) }}</td>
            <td class="px-6 py-4">{{ getWarehouseName(item.warehouseId) }}</td>
            <td class="px-6 py-4">{{ item.qtyOnHand }}</td>
            <td class="px-6 py-4">{{ item.qtyReserved }}</td>
            <td class="px-6 py-4">{{ item.qtyOnHand - item.qtyReserved }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface InventoryItem {
  id: number
  productId: number
  warehouseId: number
  qtyOnHand: number
  qtyReserved: number
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
const inventory = ref<InventoryItem[]>([])
const products = ref<Product[]>([])
const warehouses = ref<Warehouse[]>([])
const productFilter = ref<number | ''>('')
const warehouseFilter = ref<number | ''>('')
const inventoryError = ref('')

const fetchInventory = async () => {
  inventoryError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const params = new URLSearchParams()
    if (productFilter.value !== '') params.set('productId', String(productFilter.value))
    if (warehouseFilter.value !== '') params.set('warehouseId', String(warehouseFilter.value))
    const url = `/api/inventory${params.toString() ? `?${params.toString()}` : ''}`
    const response = await fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    inventory.value = result.data || []
  } catch {
    inventoryError.value = 'Gagal memuat inventaris'
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
  await Promise.all([fetchInventory(), fetchProducts(), fetchWarehouses()])
})
</script>

<style scoped>
</style>
