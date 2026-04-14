<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen Produk</h2>
        <p class="text-sm text-slate-500">Kelola SKU, harga, kategori, dan status produk.</p>
      </div>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
        <input v-model="productSearch" type="text" placeholder="Cari produk..." class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
        <button @click="openProductForm()" class="rounded-full bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Tambah Produk</button>
      </div>
    </div>

    <div v-if="productError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ productError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">SKU</th>
            <th class="px-6 py-4">Nama</th>
            <th class="px-6 py-4">Kategori</th>
            <th class="px-6 py-4">Unit</th>
            <th class="px-6 py-4">Harga</th>
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="product in filteredProducts" :key="product.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ product.sku }}</td>
            <td class="px-6 py-4">{{ product.name }}</td>
            <td class="px-6 py-4">{{ getCategoryName(product.categoryId) }}</td>
            <td class="px-6 py-4">{{ getUnitName(product.unitId) }}</td>
            <td class="px-6 py-4">Rp {{ formatNumber(product.standardPrice) }}</td>
            <td class="px-6 py-4">{{ product.lifecycleStatus }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="openProductForm(product)" class="mr-2 rounded-full bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 hover:bg-slate-200 transition">Edit</button>
              <button @click="deleteProduct(product.id)" class="rounded-full bg-red-50 px-3 py-2 text-xs font-semibold text-red-700 hover:bg-red-100 transition">Hapus</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showProductForm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeProductForm">
      <div class="w-full max-w-2xl rounded-3xl bg-white shadow-2xl ring-1 ring-slate-200 overflow-hidden max-h-[calc(100vh-3rem)] flex flex-col">
        <div class="border-b border-slate-200 p-6 shrink-0">
          <h3 class="text-xl font-semibold text-slate-900">{{ editingProduct ? 'Edit Produk' : 'Tambah Produk' }}</h3>
        </div>
        <div class="px-6 pt-4">
          <div v-if="productFormError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-4">{{ productFormError }}</div>
          <div v-if="productFormSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-4">{{ productFormSuccess }}</div>
        </div>
        <div class="p-6 grid gap-4 sm:grid-cols-2 overflow-y-auto min-h-0">
          <label class="space-y-2">
            <span class="text-sm text-slate-600">SKU</span>
            <input v-model="productForm.sku" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productFormErrors.sku" class="text-xs text-red-600">{{ productFormErrors.sku }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Nama Produk</span>
            <input v-model="productForm.name" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productFormErrors.name" class="text-xs text-red-600">{{ productFormErrors.name }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Kategori</span>
            <select v-model="productForm.categoryId" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option :value="null">Pilih kategori</option>
              <option v-for="category in categories" :key="category.id" :value="category.id">{{ category.name }}</option>
            </select>
            <p v-if="productFormErrors.categoryId" class="text-xs text-red-600">{{ productFormErrors.categoryId }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Unit</span>
            <select v-model="productForm.unitId" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option :value="null">Pilih unit</option>
              <option v-for="unit in units" :key="unit.id" :value="unit.id">{{ unit.name }}</option>
            </select>
            <p v-if="productFormErrors.unitId" class="text-xs text-red-600">{{ productFormErrors.unitId }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Harga Jual</span>
            <input v-model.number="productForm.standardPrice" type="number" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productFormErrors.standardPrice" class="text-xs text-red-600">{{ productFormErrors.standardPrice }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">HPP</span>
            <input v-model.number="productForm.costPrice" type="number" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="productFormErrors.costPrice" class="text-xs text-red-600">{{ productFormErrors.costPrice }}</p>
          </label>
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Deskripsi</span>
            <textarea v-model="productForm.description" rows="4" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none"></textarea>
          </label>
        </div>
        <div class="flex justify-end gap-3 border-t border-slate-200 p-5">
          <button @click="closeProductForm" class="rounded-2xl border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-100 transition">Batal</button>
          <button @click="saveProduct" class="rounded-2xl bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

interface Product {
  id: number
  sku: string
  name: string
  description: string
  categoryId?: number | null
  unitId?: number | null
  standardPrice: number
  costPrice: number
  lifecycleStatus: string
}

interface Category {
  id: number
  name: string
}

interface Unit {
  id: number
  name: string
}

const router = useRouter()
const products = ref<Product[]>([])
const categories = ref<Category[]>([])
const units = ref<Unit[]>([])
const productSearch = ref('')
const productError = ref('')
const productFormError = ref('')
const productFormSuccess = ref('')
const productFormErrors = ref<Record<string, string>>({})
const showProductForm = ref(false)
const editingProduct = ref(false)
const productForm = ref<Product>({ id: 0, sku: '', name: '', description: '', categoryId: null, unitId: null, standardPrice: 0, costPrice: 0, lifecycleStatus: 'active' })

const filteredProducts = computed(() => {
  const query = productSearch.value.trim().toLowerCase()
  if (!query) return products.value
  return products.value.filter((product) => product.name.toLowerCase().includes(query) || product.sku.toLowerCase().includes(query))
})

const formatNumber = (value: number | string) => {
  const num = typeof value === 'string' ? Number(value) : value
  return num.toLocaleString('id-ID')
}

const getCategoryName = (categoryId?: number | null) => categories.value.find((item) => item.id === categoryId)?.name || 'Umum'
const getUnitName = (unitId?: number | null) => units.value.find((item) => item.id === unitId)?.name || 'Tidak ditentukan'

const fetchProducts = async () => {
  productError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/products', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    products.value = result.data || []
  } catch {
    productError.value = 'Gagal memuat produk'
  }
}

const fetchCategories = async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/categories', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    categories.value = result.data || []
  } catch {
    // ignore
  }
}

const fetchUnits = async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/units', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    units.value = result.data || []
  } catch {
    // ignore
  }
}

const openProductForm = (product?: Product) => {
  productFormError.value = ''
  productFormSuccess.value = ''
  productFormErrors.value = {}

  if (product) {
    editingProduct.value = true
    productForm.value = { ...product }
  } else {
    editingProduct.value = false
    productForm.value = { id: 0, sku: '', name: '', description: '', categoryId: null, unitId: null, standardPrice: 0, costPrice: 0, lifecycleStatus: 'active' }
  }
  showProductForm.value = true
}

const closeProductForm = () => {
  showProductForm.value = false
  productFormError.value = ''
  productFormSuccess.value = ''
  productFormErrors.value = {}
}

const validateProductForm = () => {
  productFormError.value = ''
  productFormSuccess.value = ''
  productFormErrors.value = {}

  if (!productForm.value.sku.trim()) productFormErrors.value.sku = 'SKU wajib diisi.'
  if (!productForm.value.name.trim()) productFormErrors.value.name = 'Nama produk wajib diisi.'
  if (!productForm.value.categoryId) productFormErrors.value.categoryId = 'Kategori wajib dipilih.'
  if (!productForm.value.unitId) productFormErrors.value.unitId = 'Unit wajib dipilih.'
  if (productForm.value.standardPrice <= 0) productFormErrors.value.standardPrice = 'Harga jual harus lebih besar dari nol.'
  if (productForm.value.costPrice <= 0) productFormErrors.value.costPrice = 'HPP harus lebih besar dari nol.'

  if (Object.keys(productFormErrors.value).length > 0) {
    productFormError.value = 'Perbaiki kesalahan di form sebelum menyimpan.'
    return false
  }
  return true
}

const saveProduct = async () => {
  if (!validateProductForm()) return

  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  const url = editingProduct.value ? `/api/products/${productForm.value.id}` : '/api/products'
  const method = editingProduct.value ? 'PUT' : 'POST'
  try {
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(productForm.value) })
    if (!response.ok) throw new Error()
    await fetchProducts()
    productFormSuccess.value = 'Produk berhasil disimpan.'
    setTimeout(() => closeProductForm(), 800)
  } catch {
    productError.value = 'Gagal menyimpan produk'
  }
}

const deleteProduct = async (id: number) => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  if (!confirm('Hapus produk ini?')) return
  try {
    const response = await fetch(`/api/products/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    await fetchProducts()
  } catch {
    productError.value = 'Gagal menghapus produk'
  }
}

onMounted(async () => {
  await Promise.all([fetchProducts(), fetchCategories(), fetchUnits()])
})
</script>

<style scoped>
</style>
