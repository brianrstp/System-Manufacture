<template>
  <main class="bg-white text-gray-800 font-sans">
    <nav :class="['fixed top-0 left-0 w-full z-50 transition-all duration-300', scrolled ? 'bg-white/95 backdrop-blur-md shadow-lg' : 'bg-transparent']">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-20">
          <a href="#home" class="flex items-center space-x-2">
            <div class="w-10 h-10 bg-primary-700 rounded-lg flex items-center justify-center">
              <i class="fas fa-industry text-white text-lg"></i>
            </div>
            <span :class="['text-xl font-bold transition-colors', scrolled ? 'text-primary-900' : 'text-white']">
              Indo<span :class="scrolled ? 'text-primary-600' : 'text-primary-200'">Manufaktur</span>
            </span>
          </a>

          <div class="hidden md:flex items-center space-x-8">
            <a v-for="item in menuItems" :key="item.name" :href="item.href" :class="['text-sm font-medium transition-colors hover:text-primary-500', scrolled ? 'text-gray-700' : 'text-white/90']">
              {{ item.name }}
            </a>
            <a href="#contact" class="bg-primary-600 hover:bg-primary-700 text-white px-6 py-2.5 rounded-full text-sm font-semibold transition-all hover:shadow-lg hover:shadow-primary-500/30">
              Hubungi Kami
            </a>
          </div>

          <button @click="mobileMenu = !mobileMenu" class="md:hidden">
            <i :class="['fas text-xl transition-colors', mobileMenu ? 'fa-times' : 'fa-bars', scrolled ? 'text-primary-900' : 'text-white']"></i>
          </button>
        </div>
      </div>

      <transition name="slide">
        <div v-if="mobileMenu" class="md:hidden bg-white shadow-xl border-t">
          <div class="px-4 py-6 space-y-4">
            <a v-for="item in menuItems" :key="item.name" :href="item.href" @click="mobileMenu = false" class="block text-gray-700 hover:text-primary-600 font-medium py-2 border-b border-gray-100">
              {{ item.name }}
            </a>
            <a href="#contact" @click="mobileMenu = false" class="block bg-primary-600 text-white text-center px-6 py-3 rounded-full font-semibold">
              Hubungi Kami
            </a>
          </div>
        </div>
      </transition>
    </nav>

    <section id="home" class="hero-bg min-h-screen flex items-center relative overflow-hidden">
      <div class="absolute inset-0 opacity-10">
        <div class="absolute top-20 left-10 w-72 h-72 bg-white rounded-full blur-3xl"></div>
        <div class="absolute bottom-20 right-10 w-96 h-96 bg-primary-300 rounded-full blur-3xl"></div>
        <div class="absolute top-1/2 left-1/2 w-64 h-64 bg-white rounded-full blur-2xl"></div>
      </div>

      <div class="absolute inset-0 opacity-5" style="background-image: radial-gradient(circle, white 1px, transparent 1px); background-size: 40px 40px;"></div>

      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-32 relative z-10">
        <div class="grid lg:grid-cols-2 gap-12 items-center">
          <div class="text-white">
            <div class="inline-flex items-center bg-white/10 backdrop-blur-sm border border-white/20 rounded-full px-4 py-2 mb-6">
              <span class="w-2 h-2 bg-green-400 rounded-full mr-2 animate-pulse"></span>
              <span class="text-sm text-primary-100">Terpercaya Sejak 1995</span>
            </div>
            <h1 class="text-4xl sm:text-5xl lg:text-6xl font-extrabold leading-tight mb-6">
              Solusi Manufaktur
              <span class="text-primary-200">Berkualitas</span>
              Tinggi
            </h1>
            <p class="text-lg text-primary-100 mb-8 max-w-xl leading-relaxed">
              Kami menyediakan layanan manufaktur terpadu dengan teknologi terkini,
              tenaga ahli berpengalaman, dan standar kualitas internasional untuk
              memenuhi kebutuhan industri Anda.
            </p>
            <div class="flex flex-col sm:flex-row gap-4">
              <a href="#contact" class="pulse-btn bg-white text-primary-800 px-8 py-4 rounded-full font-bold text-center hover:bg-primary-50 transition-all hover:shadow-xl">
                <i class="fas fa-paper-plane mr-2"></i>Mulai Kerjasama
              </a>
              <a href="#services" class="border-2 border-white/40 text-white px-8 py-4 rounded-full font-bold text-center hover:bg-white/10 transition-all">
                <i class="fas fa-play-circle mr-2"></i>Lihat Layanan
              </a>
              <router-link to="/customer/login" class="border-2 border-white/40 text-white px-8 py-4 rounded-full font-bold text-center hover:bg-white/10 transition-all">
                <i class="fas fa-user-check mr-2"></i>Login Customer
              </router-link>
            </div>
          </div>

          <div class="hidden lg:block">
            <div class="relative">
              <div class="float-animation bg-white/10 backdrop-blur-sm border border-white/20 rounded-3xl p-8">
                <div class="grid grid-cols-2 gap-4">
                  <div class="bg-white/10 rounded-2xl p-6 text-center">
                    <i class="fas fa-cogs text-4xl text-primary-200 mb-3"></i>
                    <p class="text-white font-semibold text-sm">CNC Machining</p>
                  </div>
                  <div class="bg-white/10 rounded-2xl p-6 text-center">
                    <i class="fas fa-robot text-4xl text-primary-200 mb-3"></i>
                    <p class="text-white font-semibold text-sm">Automasi</p>
                  </div>
                  <div class="bg-white/10 rounded-2xl p-6 text-center">
                    <i class="fas fa-microchip text-4xl text-primary-200 mb-3"></i>
                    <p class="text-white font-semibold text-sm">Elektronik</p>
                  </div>
                  <div class="bg-white/10 rounded-2xl p-6 text-center">
                    <i class="fas fa-warehouse text-4xl text-primary-200 mb-3"></i>
                    <p class="text-white font-semibold text-sm">Pabrikasi</p>
                  </div>
                </div>
              </div>
              <div class="absolute -bottom-6 -left-6 bg-white rounded-2xl shadow-2xl p-4 flex items-center space-x-3">
                <div class="w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                  <i class="fas fa-check-circle text-green-600 text-xl"></i>
                </div>
                <div>
                  <p class="text-2xl font-bold text-gray-900">ISO 9001</p>
                  <p class="text-xs text-gray-500">Certified</p>
                </div>
              </div>
              <div class="absolute -top-6 -right-6 bg-white rounded-2xl shadow-2xl p-4 flex items-center space-x-3">
                <div class="w-12 h-12 bg-primary-100 rounded-full flex items-center justify-center">
                  <i class="fas fa-trophy text-primary-600 text-xl"></i>
                </div>
                <div>
                  <p class="text-2xl font-bold text-gray-900">28+</p>
                  <p class="text-xs text-gray-500">Tahun Pengalaman</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="absolute bottom-8 left-1/2 transform -translate-x-1/2 text-white/60 animate-bounce">
        <i class="fas fa-chevron-down text-2xl"></i>
      </div>
    </section>

    <section class="bg-white shadow-xl relative z-20 -mt-1 mx-4 sm:mx-8 lg:mx-16 rounded-2xl">
      <div class="grid grid-cols-2 md:grid-cols-4 divide-x divide-gray-100">
        <div v-for="(stat, index) in stats" :key="index" class="p-6 sm:p-8 text-center">
          <div class="text-3xl sm:text-4xl font-extrabold text-primary-700 mb-1">{{ stat.number }}</div>
          <div class="text-sm text-gray-500 font-medium">{{ stat.label }}</div>
        </div>
      </div>
    </section>

    <section id="about" class="py-20 sm:py-28 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid lg:grid-cols-2 gap-16 items-center">
          <div class="relative fade-in-section">
            <div class="bg-gradient-to-br from-primary-100 to-primary-200 rounded-3xl p-8 sm:p-12 relative">
              <div class="bg-white rounded-2xl shadow-xl overflow-hidden">
                <div class="bg-primary-700 p-8 text-center">
                  <i class="fas fa-industry text-6xl text-white/90 mb-4"></i>
                  <h3 class="text-white text-xl font-bold">Pabrik Modern</h3>
                  <p class="text-primary-200 text-sm mt-2">Fasilitas berstandar internasional</p>
                </div>
                <div class="p-6 grid grid-cols-3 gap-4">
                  <div class="text-center p-3 bg-primary-50 rounded-xl">
                    <i class="fas fa-shield-alt text-primary-600 text-xl mb-1"></i>
                    <p class="text-xs text-gray-600 font-medium mt-1">Kualitas Terjamin</p>
                  </div>
                  <div class="text-center p-3 bg-primary-50 rounded-xl">
                    <i class="fas fa-clock text-primary-600 text-xl mb-1"></i>
                    <p class="text-xs text-gray-600 font-medium mt-1">Tepat Waktu</p>
                  </div>
                  <div class="text-center p-3 bg-primary-50 rounded-xl">
                    <i class="fas fa-hand-holding-usd text-primary-600 text-xl mb-1"></i>
                    <p class="text-xs text-gray-600 font-medium mt-1">Harga Kompetitif</p>
                  </div>
                </div>
              </div>
            </div>
            <div class="absolute -bottom-4 -right-4 w-24 h-24 bg-primary-600 rounded-2xl -z-10"></div>
            <div class="absolute -top-4 -left-4 w-16 h-16 bg-primary-200 rounded-xl -z-10"></div>
          </div>

          <div class="fade-in-section">
            <span class="inline-block bg-primary-100 text-primary-700 text-sm font-semibold px-4 py-1.5 rounded-full mb-4">
              <i class="fas fa-info-circle mr-1"></i> Tentang Kami
            </span>
            <h2 class="text-3xl sm:text-4xl font-extrabold text-gray-900 mb-6 leading-tight">
              Mitra Manufaktur <span class="text-primary-600">Terdepan</span> di Indonesia
            </h2>
            <p class="text-gray-600 text-lg leading-relaxed mb-6">
              PT IndoManufaktur Indonesia adalah perusahaan manufaktur terkemuka yang telah
              beroperasi sejak tahun 1995. Kami berkomitmen untuk memberikan produk dan layanan
              terbaik dengan mengutamakan kualitas, efisiensi, dan inovasi.
            </p>
            <div class="space-y-4 mb-8">
              <div v-for="(item, i) in aboutFeatures" :key="i" class="flex items-start space-x-3">
                <div class="w-6 h-6 bg-primary-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                  <i class="fas fa-check text-primary-600 text-xs"></i>
                </div>
                <p class="text-gray-700">{{ item }}</p>
              </div>
            </div>
            <a href="#services" class="inline-flex items-center text-primary-600 font-semibold hover:text-primary-800 transition-colors group">
              Pelajari Lebih Lanjut
              <i class="fas fa-arrow-right ml-2 transform group-hover:translate-x-1 transition-transform"></i>
            </a>
          </div>
        </div>
      </div>
    </section>

    <section id="services" class="py-20 sm:py-28 bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center max-w-3xl mx-auto mb-16 fade-in-section">
          <span class="inline-block bg-primary-100 text-primary-700 text-sm font-semibold px-4 py-1.5 rounded-full mb-4">
            <i class="fas fa-concierge-bell mr-1"></i> Layanan Kami
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-gray-900 mb-4">
            Layanan Manufaktur <span class="text-primary-600">Komprehensif</span>
          </h2>
          <p class="text-gray-600 text-lg">
            Kami menyediakan berbagai layanan manufaktur yang dirancang untuk memenuhi
            kebutuhan industri modern dengan standar kualitas tertinggi.
          </p>
        </div>

        <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
          <div v-for="(service, index) in services" :key="index" class="fade-in-section bg-white rounded-2xl p-8 shadow-sm hover:shadow-xl transition-all duration-300 border border-gray-100 hover:border-primary-200 group cursor-pointer">
            <div class="w-16 h-16 bg-primary-50 group-hover:bg-primary-600 rounded-2xl flex items-center justify-center mb-6 transition-colors duration-300">
              <i :class="[service.icon, 'text-2xl text-primary-600 group-hover:text-white transition-colors']"></i>
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3 group-hover:text-primary-700 transition-colors">{{ service.title }}</h3>
            <p class="text-gray-600 leading-relaxed mb-4">{{ service.description }}</p>
            <div class="flex items-center text-primary-600 font-semibold text-sm group-hover:text-primary-800">
              <span>Selengkapnya</span>
              <i class="fas fa-arrow-right ml-2 transform group-hover:translate-x-1 transition-transform"></i>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="py-20 sm:py-28 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center max-w-3xl mx-auto mb-16 fade-in-section">
          <span class="inline-block bg-primary-100 text-primary-700 text-sm font-semibold px-4 py-1.5 rounded-full mb-4">
            <i class="fas fa-project-diagram mr-1"></i> Proses Kerja
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-gray-900 mb-4">
            Alur Kerja <span class="text-primary-600">Profesional</span>
          </h2>
          <p class="text-gray-600 text-lg">Proses kerja kami yang terstruktur memastikan hasil terbaik untuk setiap proyek.</p>
        </div>

        <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-8">
          <div v-for="(step, index) in processSteps" :key="index" class="fade-in-section text-center relative group">
            <div v-if="index < processSteps.length - 1" class="hidden lg:block absolute top-10 left-1/2 w-full h-0.5 bg-primary-200"></div>
            <div class="relative z-10">
              <div class="w-20 h-20 bg-primary-600 rounded-full flex items-center justify-center mx-auto mb-6 text-white text-2xl font-bold shadow-lg shadow-primary-500/30 group-hover:scale-110 transition-transform">
                {{ step.step }}
              </div>
              <h3 class="text-lg font-bold text-gray-900 mb-2">{{ step.title }}</h3>
              <p class="text-gray-600 text-sm leading-relaxed">{{ step.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="products" class="py-20 sm:py-28 hero-bg relative overflow-hidden">
      <div class="absolute inset-0 opacity-5" style="background-image: radial-gradient(circle, white 1px, transparent 1px); background-size: 30px 30px;"></div>
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div class="text-center max-w-3xl mx-auto mb-16 fade-in-section">
          <span class="inline-block bg-white/10 text-primary-200 text-sm font-semibold px-4 py-1.5 rounded-full mb-4 border border-white/20">
            <i class="fas fa-boxes mr-1"></i> Produk Kami
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-white mb-4">Produk Unggulan Kami</h2>
          <p class="text-primary-200 text-lg">Berbagai produk manufaktur berkualitas tinggi yang telah teruji dan dipercaya oleh industri nasional maupun internasional.</p>
        </div>

        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
          <div v-for="(product, index) in products" :key="index" class="fade-in-section bg-white/10 backdrop-blur-sm border border-white/20 rounded-2xl p-6 hover:bg-white/20 transition-all duration-300 group cursor-pointer">
            <div class="w-14 h-14 bg-white/10 rounded-xl flex items-center justify-center mb-4 group-hover:bg-white/20 transition-colors">
              <i :class="[product.icon, 'text-2xl text-primary-200']"></i>
            </div>
            <h3 class="text-white font-bold text-lg mb-2">{{ product.name }}</h3>
            <p class="text-primary-200 text-sm leading-relaxed">{{ product.desc }}</p>
            <div class="mt-4 pt-4 border-t border-white/10">
              <span class="text-white/60 text-xs font-medium">{{ product.category }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="py-20 sm:py-28 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid lg:grid-cols-2 gap-16 items-center">
          <div class="fade-in-section">
            <span class="inline-block bg-primary-100 text-primary-700 text-sm font-semibold px-4 py-1.5 rounded-full mb-4">
              <i class="fas fa-star mr-1"></i> Mengapa Memilih Kami
            </span>
            <h2 class="text-3xl sm:text-4xl font-extrabold text-gray-900 mb-6 leading-tight">Keunggulan yang <span class="text-primary-600">Membedakan</span> Kami</h2>
            <p class="text-gray-600 text-lg mb-8">Dengan pengalaman lebih dari 28 tahun, kami telah membuktikan diri sebagai mitra manufaktur yang dapat diandalkan.</p>
          </div>
          <div class="grid sm:grid-cols-2 gap-6">
            <div v-for="(advantage, index) in advantages" :key="index" class="fade-in-section bg-gray-50 rounded-2xl p-6 hover:bg-primary-50 transition-colors duration-300 border border-gray-100 hover:border-primary-200">
              <div class="w-12 h-12 bg-primary-100 rounded-xl flex items-center justify-center mb-4">
                <i :class="[advantage.icon, 'text-primary-600 text-xl']"></i>
              </div>
              <h3 class="font-bold text-gray-900 mb-2">{{ advantage.title }}</h3>
              <p class="text-gray-600 text-sm leading-relaxed">{{ advantage.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section class="py-20 sm:py-28 bg-gray-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center max-w-3xl mx-auto mb-16 fade-in-section">
          <span class="inline-block bg-primary-100 text-primary-700 text-sm font-semibold px-4 py-1.5 rounded-full mb-4">
            <i class="fas fa-quote-left mr-1"></i> Testimoni
          </span>
          <h2 class="text-3xl sm:text-4xl font-extrabold text-gray-900 mb-4">Apa Kata <span class="text-primary-600">Klien</span> Kami</h2>
        </div>
        <div class="grid md:grid-cols-3 gap-8">
          <div v-for="(testimonial, index) in testimonials" :key="index" class="fade-in-section bg-white rounded-2xl p-8 shadow-sm border border-gray-100 hover:shadow-lg transition-shadow">
            <div class="flex items-center mb-4">
              <i v-for="s in 5" :key="s" class="fas fa-star text-yellow-400 text-sm"></i>
            </div>
            <p class="text-gray-600 leading-relaxed mb-6 italic">"{{ testimonial.text }}"</p>
            <div class="flex items-center space-x-4 pt-4 border-t border-gray-100">
              <div class="w-12 h-12 bg-primary-100 rounded-full flex items-center justify-center">
                <span class="text-primary-700 font-bold">{{ testimonial.initials }}</span>
              </div>
              <div>
                <p class="font-bold text-gray-900">{{ testimonial.name }}</p>
                <p class="text-sm text-gray-500">{{ testimonial.position }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="contact" class="py-20 sm:py-28 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="fade-in-section hero-bg rounded-3xl p-8 sm:p-12 lg:p-16 relative overflow-hidden">
          <div class="absolute top-0 right-0 w-64 h-64 bg-white/5 rounded-full -mr-32 -mt-32"></div>
          <div class="absolute bottom-0 left-0 w-48 h-48 bg-white/5 rounded-full -ml-24 -mb-24"></div>
          <div class="grid lg:grid-cols-2 gap-12 items-center relative z-10">
            <div>
              <h2 class="text-3xl sm:text-4xl font-extrabold text-white mb-4 leading-tight">Siap Memulai Proyek Bersama Kami?</h2>
              <p class="text-primary-200 text-lg mb-8 leading-relaxed">Hubungi tim kami sekarang untuk konsultasi gratis dan dapatkan penawaran terbaik untuk kebutuhan manufaktur Anda.</p>
              <div class="flex flex-col sm:flex-row gap-4">
                <a href="tel:+622112345678" class="inline-flex items-center justify-center bg-white text-primary-800 px-8 py-4 rounded-full font-bold hover:bg-primary-50 transition-colors">
                  <i class="fas fa-phone mr-2"></i>+62 21 1234 5678
                </a>
                <a href="mailto:info@indomanufaktur.co.id" class="inline-flex items-center justify-center border-2 border-white/40 text-white px-8 py-4 rounded-full font-bold hover:bg-white/10 transition-colors">
                  <i class="fas fa-envelope mr-2"></i>Email Kami
                </a>
              </div>
            </div>

            <div class="bg-white/10 backdrop-blur-sm border border-white/20 rounded-2xl p-8">
              <h3 class="text-white text-xl font-bold mb-6">Kirim Pesan</h3>
              <form @submit.prevent="submitForm" class="space-y-4">
                <div>
                  <input v-model="form.name" type="text" placeholder="Nama Lengkap" required class="w-full bg-white/10 border border-white/20 rounded-xl px-4 py-3 text-white placeholder-white/50 focus:outline-none focus:border-white/50 focus:bg-white/15 transition-colors" />
                </div>
                <div>
                  <input v-model="form.email" type="email" placeholder="Email" required class="w-full bg-white/10 border border-white/20 rounded-xl px-4 py-3 text-white placeholder-white/50 focus:outline-none focus:border-white/50 focus:bg-white/15 transition-colors" />
                </div>
                <div>
                  <input v-model="form.phone" type="tel" placeholder="No. Telepon" class="w-full bg-white/10 border border-white/20 rounded-xl px-4 py-3 text-white placeholder-white/50 focus:outline-none focus:border-white/50 focus:bg-white/15 transition-colors" />
                </div>
                <div>
                  <textarea v-model="form.message" rows="3" placeholder="Pesan Anda..." required class="w-full bg-white/10 border border-white/20 rounded-xl px-4 py-3 text-white placeholder-white/50 focus:outline-none focus:border-white/50 focus:bg-white/15 transition-colors resize-none"></textarea>
                </div>
                <button type="submit" class="w-full bg-white text-primary-800 font-bold py-3.5 rounded-xl hover:bg-primary-50 transition-colors flex items-center justify-center">
                  <i class="fas fa-paper-plane mr-2"></i>
                  {{ formSubmitted ? 'Terkirim ✓' : 'Kirim Pesan' }}
                </button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </section>

    <footer class="bg-gray-900 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-12">
          <div class="lg:col-span-1">
            <div class="flex items-center space-x-2 mb-4">
              <div class="w-10 h-10 bg-primary-600 rounded-lg flex items-center justify-center">
                <i class="fas fa-industry text-white text-lg"></i>
              </div>
              <span class="text-xl font-bold">Indo<span class="text-primary-400">Manufaktur</span></span>
            </div>
            <p class="text-gray-400 text-sm leading-relaxed mb-6">Perusahaan manufaktur terdepan di Indonesia dengan pengalaman lebih dari 28 tahun melayani berbagai industri.</p>
            <div class="flex space-x-3">
              <a v-for="social in socialLinks" :key="social.name" :href="social.url" class="w-10 h-10 bg-gray-800 hover:bg-primary-600 rounded-lg flex items-center justify-center transition-colors">
                <i :class="[social.icon, 'text-gray-400 hover:text-white']"></i>
              </a>
            </div>
          </div>
          <div>
            <h4 class="font-bold text-lg mb-4">Navigasi</h4>
            <ul class="space-y-3">
              <li v-for="item in menuItems" :key="item.name">
                <a :href="item.href" class="text-gray-400 hover:text-primary-400 text-sm transition-colors">{{ item.name }}</a>
              </li>
            </ul>
          </div>
          <div>
            <h4 class="font-bold text-lg mb-4">Layanan</h4>
            <ul class="space-y-3">
              <li v-for="service in services.slice(0, 5)" :key="service.title">
                <a href="#services" class="text-gray-400 hover:text-primary-400 text-sm transition-colors">{{ service.title }}</a>
              </li>
            </ul>
          </div>
          <div>
            <h4 class="font-bold text-lg mb-4">Kontak</h4>
            <ul class="space-y-4">
              <li class="flex items-start space-x-3">
                <i class="fas fa-map-marker-alt text-primary-400 mt-1"></i>
                <span class="text-gray-400 text-sm">Jl. Industri Raya No. 123, Kawasan Industri Jababeka, Cikarang, Bekasi</span>
              </li>
              <li class="flex items-center space-x-3">
                <i class="fas fa-phone text-primary-400"></i>
                <span class="text-gray-400 text-sm">+62 21 1234 5678</span>
              </li>
              <li class="flex items-center space-x-3">
                <i class="fas fa-envelope text-primary-400"></i>
                <span class="text-gray-400 text-sm">info@indomanufaktur.co.id</span>
              </li>
              <li class="flex items-center space-x-3">
                <i class="fas fa-clock text-primary-400"></i>
                <span class="text-gray-400 text-sm">Sen - Jum: 08.00 - 17.00</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
      <div class="border-t border-gray-800">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 flex flex-col sm:flex-row items-center justify-between">
          <p class="text-gray-500 text-sm">&copy; 2026 IndoManufaktur. All rights reserved.</p>
          <div class="flex space-x-6 mt-4 sm:mt-0">
            <a href="#" class="text-gray-500 hover:text-gray-300 text-sm transition-colors">Privacy Policy</a>
            <a href="#" class="text-gray-500 hover:text-gray-300 text-sm transition-colors">Terms of Service</a>
          </div>
        </div>
      </div>
    </footer>

    <transition name="fade">
      <button v-if="scrolled" @click="scrollToTop" class="fixed bottom-6 right-6 w-12 h-12 bg-primary-600 text-white rounded-full shadow-lg hover:bg-primary-700 transition-all z-50 flex items-center justify-center hover:scale-110">
        <i class="fas fa-arrow-up"></i>
      </button>
    </transition>
  </main>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const scrolled = ref(false)
const mobileMenu = ref(false)
const formSubmitted = ref(false)
const form = ref({ name: '', email: '', phone: '', message: '' })

const menuItems = [
  { name: 'Beranda', href: '#home' },
  { name: 'Tentang', href: '#about' },
  { name: 'Layanan', href: '#services' },
  { name: 'Produk', href: '#products' },
  { name: 'Kontak', href: '#contact' },
]

const stats = [
  { number: '28+', label: 'Tahun Pengalaman' },
  { number: '500+', label: 'Proyek Selesai' },
  { number: '150+', label: 'Klien Aktif' },
  { number: '99.5%', label: 'Tingkat Kepuasan' },
]

const aboutFeatures = [
  'Fasilitas manufaktur modern seluas 15.000 m²',
  'Dilengkapi dengan mesin CNC dan robotika terkini',
  'Tim insinyur berpengalaman lebih dari 200 orang',
  'Sertifikasi ISO 9001:2015 dan ISO 14001:2015',
  'Melayani pasar domestik dan ekspor internasional',
  'Sistem quality control ketat di setiap tahap produksi',
]

const services = [
  {
    icon: 'fas fa-cogs',
    title: 'CNC Machining',
    description: 'Layanan pemesinan presisi tinggi dengan mesin CNC terbaru untuk berbagai material logam dan non-logam.',
  },
  {
    icon: 'fas fa-robot',
    title: 'Automasi Industri',
    description: 'Solusi automasi pabrik yang terintegrasi untuk meningkatkan efisiensi dan produktivitas lini produksi Anda.',
  },
  {
    icon: 'fas fa-layer-group',
    title: 'Sheet Metal Fabrication',
    description: 'Fabrikasi lembaran logam dengan teknik cutting, bending, dan welding berstandar industri.',
  },
  {
    icon: 'fas fa-microchip',
    title: 'PCB Assembly',
    description: 'Perakitan PCB dengan teknologi SMT dan through-hole untuk produk elektronik berkualitas tinggi.',
  },
  {
    icon: 'fas fa-drafting-compass',
    title: 'Desain & Engineering',
    description: 'Layanan desain produk dan rekayasa teknik dari konsep hingga prototipe siap produksi.',
  },
  {
    icon: 'fas fa-truck-loading',
    title: 'Supply Chain Management',
    description: 'Pengelolaan rantai pasok terpadu dari pengadaan bahan baku hingga distribusi produk jadi.',
  },
]

const processSteps = [
  { step: '01', title: 'Konsultasi', description: 'Diskusi kebutuhan dan spesifikasi proyek Anda bersama tim ahli kami.' },
  { step: '02', title: 'Desain & R&D', description: 'Tim engineering merancang solusi optimal berdasarkan kebutuhan.' },
  { step: '03', title: 'Produksi', description: 'Proses manufaktur dengan kontrol kualitas ketat di setiap tahap.' },
  { step: '04', title: 'Pengiriman', description: 'Produk dikemas aman dan dikirim tepat waktu ke lokasi Anda.' },
]

const products = [
  { icon: 'fas fa-car', name: 'Komponen Otomotif', desc: 'Suku cadang dan komponen presisi untuk industri otomotif.', category: 'Otomotif' },
  { icon: 'fas fa-bolt', name: 'Panel Elektrik', desc: 'Panel kontrol dan distribusi listrik berstandar industri.', category: 'Elektrikal' },
  { icon: 'fas fa-tools', name: 'Alat Industri', desc: 'Peralatan dan perkakas industri custom sesuai kebutuhan.', category: 'Peralatan' },
  { icon: 'fas fa-building', name: 'Struktur Baja', desc: 'Konstruksi dan fabrikasi struktur baja untuk berbagai proyek.', category: 'Konstruksi' },
]

const advantages = [
  { icon: 'fas fa-medal', title: 'Kualitas Terjamin', description: 'Standar ISO internasional dengan quality control ketat.' },
  { icon: 'fas fa-shipping-fast', title: 'Pengiriman Tepat Waktu', description: 'Track record pengiriman on-time rate di atas 98%.' },
  { icon: 'fas fa-tags', title: 'Harga Kompetitif', description: 'Penawaran harga terbaik tanpa mengorbankan kualitas.' },
  { icon: 'fas fa-headset', title: 'Support 24/7', description: 'Tim support siap membantu kapan saja Anda butuhkan.' },
]

const testimonials = [
  {
    text: 'IndoManufaktur selalu memberikan kualitas produk yang konsisten dan pengiriman yang tepat waktu. Mitra yang sangat bisa diandalkan.',
    name: 'Budi Santoso',
    position: 'Direktur, PT Otomotif Jaya',
    initials: 'BS',
  },
  {
    text: 'Kami sangat puas dengan layanan engineering mereka. Tim yang profesional dan responsif terhadap setiap kebutuhan kami.',
    name: 'Sari Dewi',
    position: 'VP Operations, TechCorp Indonesia',
    initials: 'SD',
  },
  {
    text: 'Harga kompetitif dengan kualitas premium. Sudah 10 tahun bekerjasama dan tidak pernah kecewa dengan hasilnya.',
    name: 'Ahmad Rizki',
    position: 'Procurement Manager, Global Manufacturing',
    initials: 'AR',
  },
]

const socialLinks = [
  { name: 'Facebook', icon: 'fab fa-facebook-f', url: '#' },
  { name: 'Instagram', icon: 'fab fa-instagram', url: '#' },
  { name: 'LinkedIn', icon: 'fab fa-linkedin-in', url: '#' },
  { name: 'YouTube', icon: 'fab fa-youtube', url: '#' },
]

const handleScroll = () => {
  scrolled.value = window.scrollY > 50
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const submitForm = () => {
  formSubmitted.value = true
  setTimeout(() => {
    form.value = { name: '', email: '', phone: '', message: '' }
    formSubmitted.value = false
  }, 3000)
}

const setupFadeInObserver = () => {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('visible')
        }
      })
    },
    { threshold: 0.1 }
  )

  document.querySelectorAll('.fade-in-section').forEach((el) => {
    el.classList.add('fade-in')
    observer.observe(el)
  })
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  setupFadeInObserver()
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped></style>
