<template>
  <div class="min-h-screen bg-slate-50 text-slate-900">
    <div v-if="sidebarOpen" @click="sidebarOpen = false" class="fixed inset-0 bg-black/40 z-40 lg:hidden"></div>

    <aside :class="['fixed top-0 left-0 h-full bg-white border-r border-slate-200 z-50 shadow-xl transition-transform duration-300 ease-in-out', sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0']" style="width: 280px;">
      <div class="h-full flex flex-col">
        <div class="p-6 border-b border-slate-200">
          <div class="text-primary-600 font-bold text-xl">Customer Portal</div>
          <p class="text-sm text-slate-500 mt-1">Akses pesanan dan akun Anda.</p>
        </div>
        <nav class="flex-1 p-6 overflow-y-auto">
          <ul class="space-y-2">
            <li v-for="item in navItems" :key="item.route">
              <button @click="navigate(item.route)" class="w-full text-left rounded-2xl px-4 py-3 transition hover:bg-slate-100" :class="item.route === currentRoute ? 'bg-slate-100 font-semibold text-slate-900' : 'text-slate-600'">
                {{ item.label }}
              </button>
            </li>
          </ul>
        </nav>
        <div class="p-6 border-t border-slate-200">
          <div class="text-sm text-slate-500 mb-3">Masuk sebagai</div>
          <div class="flex items-center justify-between gap-3">
            <div>
              <div class="font-semibold text-slate-900">{{ customerName }}</div>
              <div class="text-xs text-slate-500">Customer</div>
            </div>
            <button @click="logout" class="rounded-full bg-red-50 text-red-700 px-4 py-2 text-xs font-semibold hover:bg-red-100 transition">Logout</button>
          </div>
        </div>
      </div>
    </aside>

    <div class="lg:pl-[280px]">
      <header class="sticky top-0 z-30 bg-white border-b border-slate-200 backdrop-blur-xl bg-white/90">
        <div class="flex items-center justify-between gap-4 px-4 py-4 sm:px-6 lg:px-8">
          <div class="flex items-center gap-4">
            <button @click="sidebarOpen = !sidebarOpen" class="lg:hidden inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-slate-200 bg-white text-slate-700 hover:border-slate-300 hover:text-slate-900 transition">
              <span class="sr-only">Toggle sidebar</span>
              <i class="fas fa-bars"></i>
            </button>
            <div>
              <p class="text-sm text-primary-600 font-semibold">Customer Portal</p>
              <h1 class="text-2xl font-bold text-slate-900">Profil Saya</h1>
            </div>
          </div>
          <div class="hidden sm:flex items-center gap-3">
            <div class="text-right">
              <div class="text-sm font-semibold text-slate-900">{{ customerName }}</div>
              <div class="text-xs text-slate-500">Akun pelanggan</div>
            </div>
            <button @click="logout" class="rounded-2xl bg-red-50 px-4 py-2 text-sm font-semibold text-red-700 hover:bg-red-100 transition">Logout</button>
          </div>
        </div>
      </header>

      <main class="px-4 py-8 sm:px-6 lg:px-10">
        <section class="rounded-3xl bg-white border border-slate-200 shadow-sm p-6">
          <div class="mb-6">
            <p class="text-sm font-medium text-slate-500">Informasi Akun</p>
            <h2 class="text-xl font-bold text-slate-900">Profil Customer</h2>
          </div>

          <div v-if="profileError" class="rounded-2xl bg-red-50 border border-red-200 p-4 text-sm text-red-700 mb-6">
            {{ profileError }}
          </div>
          <div v-if="profileSuccess" class="rounded-2xl bg-emerald-50 border border-emerald-200 p-4 text-sm text-emerald-700 mb-6">
            {{ profileSuccess }}
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <label class="rounded-3xl border border-slate-200 bg-slate-50 p-5 block">
              <span class="text-sm text-slate-500">Nama</span>
              <input v-model="profileForm.name" type="text" class="mt-2 w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
            </label>
            <label class="rounded-3xl border border-slate-200 bg-slate-50 p-5 block">
              <span class="text-sm text-slate-500">Email</span>
              <input v-model="profileForm.email" type="email" class="mt-2 w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
            </label>
            <label class="rounded-3xl border border-slate-200 bg-slate-50 p-5 block">
              <span class="text-sm text-slate-500">Telepon</span>
              <input v-model="profileForm.phone" type="text" class="mt-2 w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
            </label>
            <label class="rounded-3xl border border-slate-200 bg-slate-50 p-5 block">
              <span class="text-sm text-slate-500">Kata Sandi Baru</span>
              <input v-model="profileForm.password" type="password" placeholder="Kosongkan jika tidak ingin mengganti" class="mt-2 w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none" />
              <p class="text-xs text-slate-500 mt-2">Minimal 8 karakter jika ingin mengganti kata sandi.</p>
            </label>
            <label class="rounded-3xl border border-slate-200 bg-slate-50 p-5 block sm:col-span-2">
              <span class="text-sm text-slate-500">Alamat</span>
              <textarea v-model="profileForm.address" rows="4" class="mt-2 w-full rounded-2xl border border-slate-200 bg-white px-4 py-3 text-sm outline-none resize-none"></textarea>
            </label>
          </div>

          <div class="mt-6 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
            <div class="text-sm text-slate-500">Email yang tersimpan: {{ customerEmail || '-' }}</div>
            <button @click.prevent="saveProfile" :disabled="isSaving" class="rounded-3xl bg-primary-600 px-6 py-3 text-sm font-semibold text-white hover:bg-primary-500 transition disabled:opacity-60">
              <span v-if="!isSaving">Simpan Profil</span>
              <span v-else>Menyimpan...</span>
            </button>
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const sidebarOpen = ref(false)
const router = useRouter()
const route = useRoute()
const customerName = ref(localStorage.getItem('customer_name') || 'Customer')
const customerEmail = ref(localStorage.getItem('customer_email') || '')
const currentRoute = computed(() => route.name as string)
const profileForm = ref({
  name: customerName.value,
  email: customerEmail.value,
  phone: '',
  address: '',
  password: '',
})
const profileError = ref('')
const profileSuccess = ref('')
const isSaving = ref(false)

const navItems = [
  { route: 'CustomerOrders', label: 'Pesanan Saya' },
  { route: 'CustomerProfile', label: 'Profil' },
  { route: 'CustomerHelp', label: 'Bantuan' },
]

const navigate = (routeName: string) => {
  router.push({ name: routeName })
}

const logout = () => {
  localStorage.removeItem('customer_token')
  localStorage.removeItem('customer_id')
  localStorage.removeItem('customer_name')
  localStorage.removeItem('customer_email')
  router.push({ name: 'CustomerLogin' })
}

const fetchProfile = async () => {
  const token = localStorage.getItem('customer_token')
  if (!token) return router.push({ name: 'CustomerLogin' })

  try {
    const response = await fetch('/api/customer/profile', { headers: { Authorization: `Bearer ${token}` } })
    if (!response.ok) {
      if (response.status === 401) {
        logout()
        return
      }
      throw new Error('Network response was not ok')
    }
    const result = await response.json()
    customerName.value = result.data.name || customerName.value
    customerEmail.value = result.data.email || customerEmail.value
    profileForm.value.name = result.data.name || profileForm.value.name
    profileForm.value.email = result.data.email || profileForm.value.email
    profileForm.value.phone = result.data.phone || ''
    profileForm.value.address = result.data.address || ''
  } catch {
    // keep stored values if fetch fails
  }
}

const isValidEmail = (email: string) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)

const saveProfile = async () => {
  profileError.value = ''
  profileSuccess.value = ''

  if (!profileForm.value.name.trim() || !profileForm.value.email.trim()) {
    profileError.value = 'Nama dan email wajib diisi.'
    return
  }
  if (!isValidEmail(profileForm.value.email.trim())) {
    profileError.value = 'Email tidak valid.'
    return
  }
  if (profileForm.value.password.trim() && profileForm.value.password.trim().length < 8) {
    profileError.value = 'Kata sandi minimal 8 karakter jika diisi.'
    return
  }

  isSaving.value = true
  try {
    const token = localStorage.getItem('customer_token')
    if (!token) return router.push({ name: 'CustomerLogin' })

    const response = await fetch('/api/customer/profile', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        name: profileForm.value.name,
        email: profileForm.value.email,
        phone: profileForm.value.phone,
        address: profileForm.value.address,
        password: profileForm.value.password,
      }),
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ message: 'Gagal memperbarui profil' }))
      profileError.value = error.message || 'Gagal memperbarui profil'
      return
    }

    const result = await response.json()
    customerName.value = result.data.name || customerName.value
    customerEmail.value = result.data.email || customerEmail.value
    localStorage.setItem('customer_name', customerName.value)
    localStorage.setItem('customer_email', customerEmail.value)
    profileForm.value.password = ''
    profileSuccess.value = 'Profil berhasil diperbarui.'
  } catch {
    profileError.value = 'Gagal memperbarui profil'
  } finally {
    isSaving.value = false
  }
}

onMounted(() => {
  fetchProfile()
})
</script>

<style scoped></style>
