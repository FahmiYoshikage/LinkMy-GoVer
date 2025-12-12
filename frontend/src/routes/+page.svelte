<script lang="ts">
	const demoLinks = [
		{ icon: 'bi-instagram', title: 'Instagram', color: '#E4405F' },
		{ icon: 'bi-youtube', title: 'YouTube', color: '#FF0000' },
		{ icon: 'bi-tiktok', title: 'TikTok', color: '#000000' },
		{ icon: 'bi-shop', title: 'Toko Online', color: '#10b981' }
	];
	
	const features = [
		{ icon: 'bi-lightning-charge-fill', title: 'Super Cepat', desc: 'Dibangun dengan Go untuk performa maksimal' },
		{ icon: 'bi-palette-fill', title: 'Tema Cantik', desc: 'Kustomisasi warna, gradien, dan layout' },
		{ icon: 'bi-graph-up', title: 'Analitik', desc: 'Lacak klik, negara, dan referrer' },
		{ icon: 'bi-layers-fill', title: 'Multi-Profil', desc: 'Buat banyak halaman link per akun' },
		{ icon: 'bi-grid-fill', title: 'Kategori', desc: 'Kelompokkan link dalam folder' },
		{ icon: 'bi-phone', title: 'Mobile First', desc: 'Tampilan sempurna di semua perangkat' }
	];
</script>

<svelte:head>
	<title>LinkMy - Satu Link. Semua Profil Anda.</title>
	<meta name="description" content="Platform link-in-bio yang powerful untuk menghubungkan audiensmu ke semua kontenmu. Gratis selamanya." />
</svelte:head>

<div class="landing">
	<!-- Header -->
	<nav class="header">
		<a href="/" class="logo">
			<i class="bi bi-link-45deg"></i>
			<span>LinkMy</span>
		</a>
		<div class="nav-links">
			<a href="#fitur">Fitur</a>
			<a href="#pricing">Harga</a>
			<a href="/login">Login</a>
			<a href="/register" class="btn btn-primary btn-sm">Daftar Gratis</a>
		</div>
	</nav>
	
	<!-- Hero Section -->
	<section class="hero">
		<div class="hero-content">
			<h1 class="hero-title">
				<em>Satu Link.</em><br/>
				<span class="text-gradient">Semua Profil Anda.</span>
			</h1>
			<p class="hero-subtitle">
				Platform link-in-bio yang powerful untuk menghubungkan audiensmu ke semua kontenmu
			</p>
			<p class="hero-desc">
				LinkMy membantu Anda mengumpulkan semua link penting - media sosial, 
				portfolio, toko online, dan lainnya - dalam satu halaman yang cantik dan mudah 
				diakses. Sempurna untuk kreator, influencer, bisnis, dan siapa saja yang ingin 
				tampil lebih profesional online.
			</p>
			
			<div class="hero-cta">
				<a href="/register" class="btn btn-primary btn-lg">
					<i class="bi bi-rocket-takeoff"></i>
					Mulai Gratis
				</a>
				<a href="#demo" class="btn btn-secondary btn-lg">
					<i class="bi bi-play-circle"></i>
					Lihat Demo
				</a>
			</div>
			
			<div class="hero-badges">
				<span><i class="bi bi-check-circle-fill"></i> Gratis selamanya</span>
				<span><i class="bi bi-check-circle-fill"></i> Tanpa iklan</span>
				<span><i class="bi bi-check-circle-fill"></i> Unlimited links</span>
			</div>
		</div>
		
		<!-- Phone Mockup -->
		<div class="phone-mockup">
			<div class="phone-frame">
				<div class="phone-notch"></div>
				<div class="phone-screen">
					<div class="demo-profile">
						<div class="demo-avatar">YO</div>
						<h3>@username</h3>
						<p>Kreator • Entrepreneur • Dreamer</p>
						
						<div class="demo-links">
							{#each demoLinks as link, i}
								<div class="demo-link" style="animation-delay: {i * 150}ms">
									<i class="bi {link.icon}" style="color: {link.color}"></i>
									<span>{link.title}</span>
								</div>
							{/each}
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>
	
	<!-- Features Section -->
	<section id="fitur" class="features">
		<h2 class="section-title">
			<span class="text-gradient">Fitur Powerful</span>
		</h2>
		<p class="section-subtitle">Semua yang kamu butuhkan untuk halaman link yang sempurna</p>
		
		<div class="features-grid">
			{#each features as feature, i}
				<div class="feature-card" style="animation-delay: {i * 100}ms">
					<div class="feature-icon">
						<i class="bi {feature.icon}"></i>
					</div>
					<h3>{feature.title}</h3>
					<p>{feature.desc}</p>
				</div>
			{/each}
		</div>
	</section>
	
	<!-- Pricing Section -->
	<section id="pricing" class="pricing">
		<h2 class="section-title">
			<span class="text-gradient">Harga Simpel</span>
		</h2>
		<p class="section-subtitle">Mulai gratis, upgrade kapanpun untuk fitur lebih</p>
		
		<div class="pricing-grid">
			<!-- Free Plan -->
			<div class="pricing-card">
				<div class="plan-header">
					<h3>Free</h3>
					<div class="plan-price">
						<span class="currency">$</span>
						<span class="amount">0</span>
						<span class="period">/bulan</span>
					</div>
					<p>Selamanya gratis</p>
				</div>
				<ul class="plan-features">
					<li><i class="bi bi-check-lg"></i> 2 Profil</li>
					<li><i class="bi bi-check-lg"></i> Unlimited Links</li>
					<li><i class="bi bi-check-lg"></i> Tema Basic</li>
					<li><i class="bi bi-check-lg"></i> Analitik 7 Hari</li>
					<li><i class="bi bi-check-lg"></i> Kategori/Folder</li>
					<li><i class="bi bi-check-lg"></i> Share & QR Code</li>
				</ul>
				<a href="/register" class="btn btn-secondary btn-full">
					Mulai Gratis
				</a>
			</div>
			
			<!-- Pro Plan -->
			<div class="pricing-card featured">
				<div class="plan-badge">Popular</div>
				<div class="plan-header">
					<h3>Pro</h3>
					<div class="plan-price">
						<span class="currency">$</span>
						<span class="amount">2</span>
						<span class="period">/bulan</span>
					</div>
					<p>Untuk yang serius</p>
				</div>
				<ul class="plan-features">
					<li><i class="bi bi-check-lg"></i> <strong>Unlimited</strong> Profil</li>
					<li><i class="bi bi-check-lg"></i> Unlimited Links</li>
					<li><i class="bi bi-check-lg"></i> <strong>Semua</strong> Tema</li>
					<li><i class="bi bi-check-lg"></i> Analitik <strong>Lifetime</strong></li>
					<li><i class="bi bi-check-lg"></i> Kategori/Folder</li>
					<li><i class="bi bi-check-lg"></i> Priority Support</li>
					<li><i class="bi bi-check-lg"></i> Custom Domain (soon)</li>
				</ul>
				<a href="/register" class="btn btn-primary btn-full">
					Upgrade ke Pro
				</a>
			</div>
		</div>
	</section>
	
	<!-- CTA Section -->
	<section class="cta">
		<div class="cta-content">
			<h2>Siap untuk memulai?</h2>
			<p>Buat halaman linkmu dalam hitungan detik</p>
			<a href="/register" class="btn btn-primary btn-lg">
				<i class="bi bi-arrow-right"></i>
				Daftar Gratis Sekarang
			</a>
		</div>
	</section>
	
	<!-- Footer -->
	<footer class="footer">
		<div class="footer-content">
			<div class="footer-brand">
				<a href="/" class="logo">
					<i class="bi bi-link-45deg"></i>
					<span>LinkMy</span>
				</a>
				<p>Platform link-in-bio terbaik Indonesia</p>
			</div>
			<div class="footer-links">
				<a href="/login">Login</a>
				<a href="/register">Daftar</a>
				<a href="#fitur">Fitur</a>
				<a href="#pricing">Harga</a>
			</div>
		</div>
		<div class="footer-bottom">
			<p>&copy; 2025 LinkMy. Made with <i class="bi bi-heart-fill"></i> by Fahmi</p>
		</div>
	</footer>
</div>

<style>
	.landing {
		min-height: 100vh;
	}
	
	/* Header */
	.header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--space-lg) var(--space-2xl);
		position: sticky;
		top: 0;
		background: rgba(13, 13, 13, 0.8);
		backdrop-filter: blur(10px);
		z-index: 100;
	}
	
	.logo {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--color-text);
	}
	
	.logo i {
		color: var(--color-primary);
		font-size: 1.75rem;
	}
	
	.nav-links {
		display: flex;
		align-items: center;
		gap: var(--space-xl);
	}
	
	.nav-links a {
		color: var(--color-text-secondary);
		font-size: 0.9rem;
	}
	
	.nav-links a:hover {
		color: var(--color-text);
	}
	
	.btn-sm {
		padding: var(--space-sm) var(--space-lg);
		font-size: 0.875rem;
	}
	
	/* Hero */
	.hero {
		display: grid;
		grid-template-columns: 1fr 1fr;
		gap: var(--space-3xl);
		align-items: center;
		padding: var(--space-3xl) var(--space-2xl);
		min-height: 80vh;
	}
	
	.hero-title {
		font-size: clamp(2.5rem, 5vw, 4rem);
		font-weight: 700;
		line-height: 1.1;
		margin-bottom: var(--space-xl);
	}
	
	.hero-title em {
		font-style: italic;
		font-weight: 700;
	}
	
	.hero-subtitle {
		font-size: 1.25rem;
		color: var(--color-text-secondary);
		margin-bottom: var(--space-lg);
	}
	
	.hero-desc {
		color: var(--color-text-muted);
		line-height: 1.7;
		margin-bottom: var(--space-2xl);
	}
	
	.hero-cta {
		display: flex;
		gap: var(--space-md);
		margin-bottom: var(--space-xl);
	}
	
	.btn-lg {
		padding: var(--space-lg) var(--space-2xl);
		font-size: 1rem;
	}
	
	.hero-badges {
		display: flex;
		gap: var(--space-lg);
		flex-wrap: wrap;
	}
	
	.hero-badges span {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
		font-size: 0.875rem;
		color: var(--color-text-muted);
	}
	
	.hero-badges i {
		color: var(--color-success);
	}
	
	/* Phone Mockup */
	.phone-mockup {
		display: flex;
		justify-content: center;
	}
	
	.phone-frame {
		width: 300px;
		background: var(--color-bg-secondary);
		border: 8px solid #333;
		border-radius: 40px;
		padding: var(--space-sm);
		box-shadow: 
			0 40px 80px rgba(0, 0, 0, 0.5),
			inset 0 2px 0 rgba(255, 255, 255, 0.1);
	}
	
	.phone-notch {
		width: 100px;
		height: 25px;
		background: #333;
		border-radius: 20px;
		margin: 0 auto var(--space-md);
	}
	
	.phone-screen {
		background: var(--color-bg);
		border-radius: 30px;
		padding: var(--space-xl);
		min-height: 450px;
	}
	
	.demo-profile {
		text-align: center;
	}
	
	.demo-avatar {
		width: 80px;
		height: 80px;
		background: var(--color-primary);
		border-radius: 50%;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.75rem;
		font-weight: 700;
		color: white;
		margin: 0 auto var(--space-md);
	}
	
	.demo-profile h3 {
		font-size: 1.25rem;
		margin-bottom: var(--space-xs);
	}
	
	.demo-profile > p {
		font-size: 0.75rem;
		color: var(--color-text-muted);
		margin-bottom: var(--space-xl);
	}
	
	.demo-links {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
	}
	
	.demo-link {
		display: flex;
		align-items: center;
		gap: var(--space-md);
		padding: var(--space-md) var(--space-lg);
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		animation: slideUp 0.5s ease backwards;
	}
	
	.demo-link i {
		font-size: 1rem;
	}
	
	.demo-link span {
		font-size: 0.875rem;
	}
	
	/* Features */
	.features {
		padding: var(--space-3xl) var(--space-2xl);
		background: var(--color-bg-secondary);
	}
	
	.section-title {
		font-size: 2.5rem;
		font-weight: 700;
		text-align: center;
		margin-bottom: var(--space-md);
	}
	
	.section-subtitle {
		text-align: center;
		color: var(--color-text-secondary);
		margin-bottom: var(--space-3xl);
	}
	
	.features-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
		gap: var(--space-xl);
		max-width: 1200px;
		margin: 0 auto;
	}
	
	.feature-card {
		text-align: center;
		padding: var(--space-2xl);
		background: var(--color-bg);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		animation: slideUp 0.5s ease backwards;
	}
	
	.feature-icon {
		width: 60px;
		height: 60px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--gradient-primary);
		border-radius: var(--radius-lg);
		margin: 0 auto var(--space-lg);
	}
	
	.feature-icon i {
		font-size: 1.5rem;
		color: white;
	}
	
	.feature-card h3 {
		font-size: 1.125rem;
		margin-bottom: var(--space-sm);
	}
	
	.feature-card p {
		color: var(--color-text-secondary);
		font-size: 0.875rem;
	}
	
	/* Pricing */
	.pricing {
		padding: var(--space-3xl) var(--space-2xl);
	}
	
	.pricing-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
		gap: var(--space-xl);
		max-width: 800px;
		margin: 0 auto;
	}
	
	.pricing-card {
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		padding: var(--space-2xl);
		position: relative;
	}
	
	.pricing-card.featured {
		border-color: var(--color-primary);
		background: linear-gradient(135deg, rgba(102, 126, 234, 0.1), rgba(118, 75, 162, 0.1));
	}
	
	.plan-badge {
		position: absolute;
		top: -12px;
		left: 50%;
		transform: translateX(-50%);
		background: var(--gradient-primary);
		color: white;
		padding: var(--space-xs) var(--space-lg);
		border-radius: var(--radius-full);
		font-size: 0.75rem;
		font-weight: 600;
	}
	
	.plan-header {
		text-align: center;
		margin-bottom: var(--space-xl);
	}
	
	.plan-header h3 {
		font-size: 1.5rem;
		margin-bottom: var(--space-md);
	}
	
	.plan-price {
		display: flex;
		align-items: flex-start;
		justify-content: center;
		gap: var(--space-xs);
		margin-bottom: var(--space-sm);
	}
	
	.plan-price .currency {
		font-size: 1.5rem;
		margin-top: var(--space-sm);
	}
	
	.plan-price .amount {
		font-size: 4rem;
		font-weight: 700;
		line-height: 1;
	}
	
	.plan-price .period {
		color: var(--color-text-muted);
		margin-top: auto;
		padding-bottom: var(--space-sm);
	}
	
	.plan-header > p {
		color: var(--color-text-muted);
		font-size: 0.875rem;
	}
	
	.plan-features {
		list-style: none;
		margin-bottom: var(--space-xl);
	}
	
	.plan-features li {
		display: flex;
		align-items: center;
		gap: var(--space-sm);
		padding: var(--space-sm) 0;
		border-bottom: 1px solid var(--color-border);
		font-size: 0.9rem;
	}
	
	.plan-features li:last-child {
		border-bottom: none;
	}
	
	.plan-features i {
		color: var(--color-success);
	}
	
	.btn-full {
		width: 100%;
		justify-content: center;
	}
	
	/* CTA */
	.cta {
		padding: var(--space-3xl) var(--space-2xl);
	}
	
	.cta-content {
		text-align: center;
		padding: var(--space-3xl);
		background: var(--gradient-primary);
		border-radius: var(--radius-2xl);
		max-width: 800px;
		margin: 0 auto;
	}
	
	.cta-content h2 {
		font-size: 2rem;
		margin-bottom: var(--space-md);
	}
	
	.cta-content p {
		margin-bottom: var(--space-xl);
		opacity: 0.9;
	}
	
	.cta-content .btn {
		background: white;
		color: var(--color-primary);
	}
	
	/* Footer */
	.footer {
		padding: var(--space-3xl) var(--space-2xl);
		border-top: 1px solid var(--color-border);
	}
	
	.footer-content {
		display: flex;
		justify-content: space-between;
		align-items: center;
		flex-wrap: wrap;
		gap: var(--space-xl);
		max-width: 1200px;
		margin: 0 auto var(--space-xl);
	}
	
	.footer-brand p {
		color: var(--color-text-muted);
		font-size: 0.875rem;
		margin-top: var(--space-sm);
	}
	
	.footer-links {
		display: flex;
		gap: var(--space-xl);
	}
	
	.footer-links a {
		color: var(--color-text-secondary);
		font-size: 0.9rem;
	}
	
	.footer-bottom {
		text-align: center;
		color: var(--color-text-muted);
		font-size: 0.875rem;
	}
	
	.footer-bottom i {
		color: var(--color-danger);
	}
	
	/* Animations */
	@keyframes slideUp {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
	
	/* Responsive */
	@media (max-width: 900px) {
		.hero {
			grid-template-columns: 1fr;
			text-align: center;
		}
		
		.hero-cta {
			justify-content: center;
		}
		
		.hero-badges {
			justify-content: center;
		}
		
		.phone-mockup {
			order: -1;
		}
		
		.phone-frame {
			width: 260px;
		}
		
		.phone-screen {
			min-height: 380px;
		}
	}
	
	@media (max-width: 600px) {
		.header {
			padding: var(--space-md);
		}
		
		.nav-links {
			gap: var(--space-md);
		}
		
		.nav-links a:not(.btn) {
			display: none;
		}
		
		.hero {
			padding: var(--space-xl);
		}
		
		.hero-cta {
			flex-direction: column;
		}
		
		.footer-content {
			flex-direction: column;
			text-align: center;
		}
	}
</style>
