<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h2 class="text-xl font-bold text-slate-900">Manajemen Pesanan</h2>
        <p class="text-sm text-slate-500">Kelola pesanan dan status pemenuhan.</p>
      </div>
      <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
        <input v-model="orderSearch" type="text" placeholder="Cari pesanan..." class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
        <select v-model="orderStatus" class="rounded-full border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
          <option value="">Semua Status</option>
          <option value="pending">Pending</option>
          <option value="Diproses">Diproses</option>
          <option value="Selesai">Selesai</option>
          <option value="Quality Check">Quality Check</option>
          <option value="Dikirim">Dikirim</option>
        </select>
        <button @click="fetchOrders" class="rounded-full bg-slate-100 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-200 transition">Refresh</button>
        <button @click="openOrderForm()" class="rounded-full bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Tambah Pesanan</button>
      </div>
    </div>

    <div v-if="orderError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700">{{ orderError }}</div>

    <div class="overflow-x-auto bg-white rounded-3xl border border-slate-200 shadow-sm">
      <table class="min-w-full text-left text-sm">
        <thead class="bg-slate-50 text-slate-500 uppercase tracking-wider text-xs">
          <tr>
            <th class="px-6 py-4">Nomor</th>
            <th class="px-6 py-4">Pelanggan</th>
            <th class="px-6 py-4">Produk</th>
            <th class="px-6 py-4">Tanggal</th>
            <th class="px-6 py-4">Jumlah</th>
            <th class="px-6 py-4">Status</th>
            <th class="px-6 py-4 text-right">Aksi</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-slate-100 bg-white">
          <tr v-for="order in filteredOrders" :key="order.id" class="hover:bg-slate-50 transition-colors">
            <td class="px-6 py-4">{{ order.orderNumber }}</td>
            <td class="px-6 py-4">{{ order.customerName }}</td>
            <td class="px-6 py-4">{{ order.product }}</td>
            <td class="px-6 py-4">{{ order.orderDate }}</td>
            <td class="px-6 py-4">{{ order.amount }}</td>
            <td class="px-6 py-4">{{ order.status }}</td>
            <td class="px-6 py-4 text-right">
              <button @click="openOrderForm(order)" class="mr-2 rounded-full bg-slate-100 px-3 py-2 text-xs font-semibold text-slate-700 hover:bg-slate-200 transition">Edit</button>
              <button @click="deleteOrder(order.id)" class="rounded-full bg-red-50 px-3 py-2 text-xs font-semibold text-red-700 hover:bg-red-100 transition">Hapus</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showOrderForm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4" @click.self="closeOrderForm">
      <div class="w-full max-w-2xl rounded-3xl bg-white shadow-2xl ring-1 ring-slate-200 overflow-hidden max-h-[calc(100vh-3rem)] flex flex-col">
        <div class="border-b border-slate-200 p-6 shrink-0">
          <h3 class="text-xl font-semibold text-slate-900">{{ editingOrder ? 'Edit Pesanan' : 'Tambah Pesanan' }}</h3>
        </div>
        <div class="px-6 pt-4">
          <div v-if="orderFormError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-4">{{ orderFormError }}</div>
          <div v-if="orderFormSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-4">{{ orderFormSuccess }}</div>
        </div>
        <div class="p-6 grid gap-4 sm:grid-cols-2 overflow-y-auto min-h-0">
          <label class="space-y-2 sm:col-span-2">
            <span class="text-sm text-slate-600">Nomor Pesanan</span>
            <input v-model="orderForm.orderNumber" type="text" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="orderFormErrors.orderNumber" class="text-xs text-red-600">{{ orderFormErrors.orderNumber }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Pelanggan</span>
            <select v-model="orderForm.customerId" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option :value="0">Pilih pelanggan</option>
              <option v-for="customer in customers" :key="customer.id" :value="customer.id">{{ customer.name }}</option>
            </select>
            <p v-if="orderFormErrors.customerId" class="text-xs text-red-600">{{ orderFormErrors.customerId }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Produk</span>
            <select v-model="orderForm.product" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option value="">Pilih produk</option>
              <option v-for="product in products" :key="product.id" :value="product.name">{{ product.name }}</option>
            </select>
            <p v-if="orderFormErrors.product" class="text-xs text-red-600">{{ orderFormErrors.product }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Tanggal Pesanan</span>
            <input v-model="orderForm.orderDate" type="date" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="orderFormErrors.orderDate" class="text-xs text-red-600">{{ orderFormErrors.orderDate }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Jumlah</span>
            <input v-model.number="orderForm.amount" type="number" class="w-full rounded-2xl border border-slate-200 bg-slate-50 px-4 py-3 text-sm outline-none" />
            <p v-if="orderFormErrors.amount" class="text-xs text-red-600">{{ orderFormErrors.amount }}</p>
          </label>
          <label class="space-y-2">
            <span class="text-sm text-slate-600">Status</span>
            <select v-model="orderForm.status" class="w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none">
              <option value="pending">Pending</option>
              <option value="Diproses">Diproses</option>
              <option value="Selesai">Selesai</option>
              <option value="Quality Check">Quality Check</option>
              <option value="Dikirim">Dikirim</option>
            </select>
          </label>
        </div>
        <div class="flex justify-end gap-3 border-t border-slate-200 p-5">
          <button @click="closeOrderForm" class="rounded-2xl border border-slate-200 px-5 py-3 text-sm font-semibold text-slate-700 hover:bg-slate-100 transition">Batal</button>
          <button @click="saveOrder" class="rounded-2xl bg-primary-600 px-5 py-3 text-sm font-semibold text-white hover:bg-primary-700 transition">Simpan</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

interface Order {
  id: number
  orderNumber: string
  customerId: number
  customerName: string
  product: string
  orderDate: string
  amount: number | string
  status: string
}

interface Customer {
  id: number
  name: string
}

interface Product {
  id: number
  name: string
}

const router = useRouter()
const orders = ref<Order[]>([])
const customers = ref<Customer[]>([])
const products = ref<Product[]>([])
const orderSearch = ref('')
const orderStatus = ref('')
const orderError = ref('')
const orderFormError = ref('')
const orderFormSuccess = ref('')
const orderFormErrors = ref<Record<string, string>>({})
const showOrderForm = ref(false)
const editingOrder = ref(false)
const orderForm = ref({ id: 0, orderNumber: '', customerId: 0, product: '', orderDate: new Date().toISOString().slice(0, 10), amount: 0, status: 'pending' })

const filteredOrders = computed(() => {
  const query = orderSearch.value.trim().toLowerCase()
  return orders.value.filter((order) => {
    const matchesSearch =
      query === '' ||
      order.orderNumber.toLowerCase().includes(query) ||
      order.customerName.toLowerCase().includes(query) ||
      order.product.toLowerCase().includes(query) ||
      order.status.toLowerCase().includes(query)
    const matchesStatus = orderStatus.value === '' || order.status.toLowerCase() === orderStatus.value.toLowerCase()
    return matchesSearch && matchesStatus
  })
})

const fetchOrders = async () => {
  orderError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const url = `/api/orders` + (orderStatus.value || orderSearch.value ? `?${new URLSearchParams({ search: orderSearch.value, status: orderStatus.value }).toString()}` : '')
    const response = await fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    orders.value = (result.data || []).map((order: any) => ({
      ...order,
      orderDate: typeof order.orderDate === 'string' ? order.orderDate.slice(0, 10) : order.orderDate,
    }))
  } catch {
    orderError.value = 'Gagal memuat pesanan'
  }
}

const fetchCustomers = async () => {
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  try {
    const response = await fetch('/api/customers', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    const result = await response.json()
    customers.value = result.data || []
  } catch {
    // ignore
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

const openOrderForm = (order?: Order) => {
  orderFormError.value = ''
  orderFormSuccess.value = ''
  orderFormErrors.value = {}

  if (order) {
    editingOrder.value = true
    orderForm.value = {
      id: order.id,
      orderNumber: order.orderNumber,
      customerId: order.customerId,
      product: order.product,
      orderDate: typeof order.orderDate === 'string' ? order.orderDate.slice(0, 10) : order.orderDate,
      amount: Number(order.amount) || 0,
      status: order.status,
    }
  } else {
    editingOrder.value = false
    orderForm.value = { id: 0, orderNumber: '', customerId: customers.value[0]?.id ?? 0, product: products.value[0]?.name ?? '', orderDate: new Date().toISOString().slice(0, 10), amount: 0, status: 'pending' }
  }
  showOrderForm.value = true
}

const closeOrderForm = () => {
  showOrderForm.value = false
  orderFormError.value = ''
  orderFormSuccess.value = ''
  orderFormErrors.value = {}
}

const validateOrderForm = () => {
  orderFormError.value = ''
  orderFormSuccess.value = ''
  orderFormErrors.value = {}

  if (!orderForm.value.orderNumber.trim()) orderFormErrors.value.orderNumber = 'Nomor pesanan wajib diisi.'
  if (orderForm.value.customerId <= 0) orderFormErrors.value.customerId = 'Pilih pelanggan.'
  if (!orderForm.value.product.trim()) orderFormErrors.value.product = 'Produk wajib dipilih.'
  if (!orderForm.value.orderDate) orderFormErrors.value.orderDate = 'Tanggal pesanan wajib diisi.'
  if (!orderForm.value.amount || orderForm.value.amount <= 0) orderFormErrors.value.amount = 'Jumlah harus lebih besar dari nol.'

  if (Object.keys(orderFormErrors.value).length > 0) {
    orderFormError.value = 'Perbaiki kesalahan di form sebelum menyimpan.'
    return false
  }
  return true
}

const saveOrder = async () => {
  if (!validateOrderForm()) return

  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  const url = editingOrder.value ? `/api/orders/${orderForm.value.id}` : '/api/orders'
  const method = editingOrder.value ? 'PUT' : 'POST'
  try {
    const response = await fetch(url, { method, headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` }, body: JSON.stringify(orderForm.value) })
    if (!response.ok) throw new Error()
    await fetchOrders()
    orderFormSuccess.value = 'Pesanan berhasil disimpan.'
    setTimeout(() => closeOrderForm(), 800)
  } catch {
    orderError.value = 'Gagal menyimpan pesanan'
  }
}

const deleteOrder = async (id: number) => {
  orderError.value = ''
  const token = localStorage.getItem('admin_token')
  if (!token) return router.push({ name: 'Login' })
  if (!confirm('Hapus pesanan ini?')) return
  try {
    const response = await fetch(`/api/orders/${id}`, { method: 'DELETE', headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) throw new Error()
    await fetchOrders()
  } catch {
    orderError.value = 'Gagal menghapus pesanan'
  }
}

onMounted(async () => {
  await Promise.all([fetchOrders(), fetchCustomers(), fetchProducts()])
})
</script>

<style scoped>
</style>
