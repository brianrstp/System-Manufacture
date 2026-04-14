<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen BOM</h2>
        <p class="text-sm text-slate-500">Kelola struktur bahan dan komponen produk.</p>
      </div>
      <button @click="fetchBOMs" class="rounded-full border border-slate-200 bg-white px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-50 transition">Segarkan</button>
    </div>

    <div v-if="bomError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ bomError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Produk</th>
            <th class="px-6 py-4">Komponen</th>
            <th class="px-6 py-4">Qty</th>
            <th class="px-6 py-4">Waste %</th>
            <th class="px-6 py-4">Parent BOM</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="bom in boms" :key="bom.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ getProductName(bom.productId) }}</td>
            <td class="px-6 py-4">{{ getProductName(bom.componentProductId) }}</td>
            <td class="px-6 py-4">{{ bom.componentQty }}</td>
            <td class="px-6 py-4">{{ bom.wastePercentage }}%</td>
            <td class="px-6 py-4">{{ bom.parentBomId ? bom.parentBomId : '-' }}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

interface BOMLine {
  id: number
  productId: number
  componentProductId: number
  componentQty: number
  wastePercentage: number
  parentBomId?: number | null
}

interface Product {
  id: number
  name: string
}

const router = useRouter()
const boms = ref<BOMLine[]>([])
const products = ref<Product[]>([])
const bomError = ref('')

const fetchBOMs = async () => {
  bomError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/boms', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    boms.value = result.data || []
  } catch {
    bomError.value = 'Gagal memuat BOM'
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

const getProductName = (productId: number) => products.value.find((item) => item.id === productId)?.name || `#${productId}`

onMounted(async () => {
  await Promise.all([fetchBOMs(), fetchProducts()])
})
</script>

<style scoped>
</style>
