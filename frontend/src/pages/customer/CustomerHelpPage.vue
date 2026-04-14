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
              <h1 class="text-2xl font-bold text-slate-900">Bantuan</h1>
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
            <p class="text-sm font-medium text-slate-500">Bantuan</p>
            <h2 class="text-xl font-bold text-slate-900">Pusat Bantuan</h2>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div class="rounded-3xl border border-slate-200 bg-slate-50 p-5">
              <p class="text-sm text-slate-500">Kontak Dukungan</p>
              <p class="mt-2 text-lg font-semibold text-slate-900">support@manufacture.local</p>
              <p class="text-sm text-slate-500 mt-1">(021) 555-1234</p>
            </div>
            <div class="rounded-3xl border border-slate-200 bg-slate-50 p-5 sm:col-span-2">
              <p class="text-sm text-slate-500">Pertanyaan Umum</p>
              <ul class="mt-3 space-y-3 text-slate-700">
                <li>
                  <p class="font-semibold">Bagaimana cara memeriksa status pesanan?</p>
                  <p class="text-sm text-slate-600">Pergi ke halaman Pesanan Saya untuk melihat daftar pesanan terbaru dan statusnya.</p>
                </li>
                <li>
                  <p class="font-semibold">Bagaimana cara memperbarui profil?</p>
                  <p class="text-sm text-slate-600">Untuk saat ini, perubahan profil dilakukan oleh admin. Hubungi dukungan jika perlu pembaruan.</p>
                </li>
                <li>
                  <p class="font-semibold">Bagaimana cara logout?</p>
                  <p class="text-sm text-slate-600">Gunakan tombol Logout di sidebar atau header untuk mengakhiri sesi.</p>
                </li>
              </ul>
            </div>
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const sidebarOpen = ref(false)
const router = useRouter()
const route = useRoute()
const customerName = localStorage.getItem('customer_name') || 'Customer'
const currentRoute = computed(() => route.name as string)

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
</script>

<style scoped></style>
