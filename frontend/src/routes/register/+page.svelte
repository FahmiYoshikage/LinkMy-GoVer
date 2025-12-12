<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/api';
	
	// Step: 1 = email/password, 2 = OTP, 3 = username/slug
	let step = $state(1);
	let loading = $state(false);
	let error = $state('');
	
	// Step 1 data
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	
	// Step 2 data
	let otp = $state('');
	let resendCooldown = $state(0);
	
	// Step 3 data
	let username = $state('');
	let slug = $state('');
	
	// Step 1: Send OTP
	async function submitStep1(e: Event) {
		e.preventDefault();
		error = '';
		
		if (password !== confirmPassword) {
			error = 'Password tidak sama';
			return;
		}
		
		if (password.length < 8) {
			error = 'Password minimal 8 karakter';
			return;
		}
		
		loading = true;
		const res = await auth.sendOTP(email, password);
		loading = false;
		
		if (res.error) {
			error = res.error;
			return;
		}
		
		// Start cooldown
		resendCooldown = 60;
		const interval = setInterval(() => {
			resendCooldown--;
			if (resendCooldown <= 0) clearInterval(interval);
		}, 1000);
		
		step = 2;
	}
	
	// Step 2: Verify OTP
	async function submitStep2(e: Event) {
		e.preventDefault();
		error = '';
		
		if (otp.length !== 6) {
			error = 'Masukkan kode 6 digit';
			return;
		}
		
		loading = true;
		const res = await auth.verifyOTP(email, otp);
		loading = false;
		
		if (res.error) {
			error = res.error;
			return;
		}
		
		step = 3;
	}
	
	// Resend OTP
	async function resendOTP() {
		if (resendCooldown > 0) return;
		
		loading = true;
		error = '';
		const res = await auth.sendOTP(email, password);
		loading = false;
		
		if (res.error) {
			error = res.error;
			return;
		}
		
		resendCooldown = 60;
		const interval = setInterval(() => {
			resendCooldown--;
			if (resendCooldown <= 0) clearInterval(interval);
		}, 1000);
	}
	
	// Step 3: Complete registration
	async function submitStep3(e: Event) {
		e.preventDefault();
		error = '';
		
		if (username.length < 3) {
			error = 'Username minimal 3 karakter';
			return;
		}
		
		loading = true;
		const res = await auth.completeRegistration(email, password, otp, username, slug || username);
		loading = false;
		
		if (res.error) {
			error = res.error;
			return;
		}
		
		goto('/dashboard');
	}
	
	// Auto-generate slug from username
	function generateSlug() {
		slug = username.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
	}
</script>

<svelte:head>
	<title>Daftar - LinkMy</title>
</svelte:head>

<div class="register-page">
	<div class="register-container">
		<div class="register-header">
			<a href="/" class="logo">
				<i class="bi bi-link-45deg"></i>
				<span>LinkMy</span>
			</a>
			<h1>Buat Akun</h1>
			<p>Langkah {step} dari 3</p>
		</div>
		
		<!-- Progress Steps -->
		<div class="steps">
			<div class="step" class:active={step >= 1} class:completed={step > 1}>
				<span class="step-num">1</span>
				<span class="step-label">Email</span>
			</div>
			<div class="step-line" class:active={step > 1}></div>
			<div class="step" class:active={step >= 2} class:completed={step > 2}>
				<span class="step-num">2</span>
				<span class="step-label">Verifikasi</span>
			</div>
			<div class="step-line" class:active={step > 2}></div>
			<div class="step" class:active={step >= 3}>
				<span class="step-num">3</span>
				<span class="step-label">Profil</span>
			</div>
		</div>
		
		{#if error}
			<div class="alert alert-error">
				<i class="bi bi-exclamation-circle"></i>
				{error}
			</div>
		{/if}
		
		<!-- Step 1: Email & Password -->
		{#if step === 1}
			<form class="form" onsubmit={submitStep1}>
				<div class="form-group">
					<label for="email">Email</label>
					<input
						type="email"
						id="email"
						bind:value={email}
						placeholder="nama@email.com"
						required
					/>
				</div>
				
				<div class="form-group">
					<label for="password">Password</label>
					<input
						type="password"
						id="password"
						bind:value={password}
						placeholder="Minimal 8 karakter"
						required
					/>
				</div>
				
				<div class="form-group">
					<label for="confirmPassword">Konfirmasi Password</label>
					<input
						type="password"
						id="confirmPassword"
						bind:value={confirmPassword}
						placeholder="Ulangi password"
						required
					/>
				</div>
				
				<button type="submit" class="btn btn-primary btn-full" disabled={loading}>
					{#if loading}
						<i class="bi bi-arrow-repeat spin"></i>
						Mengirim OTP...
					{:else}
						Lanjut
						<i class="bi bi-arrow-right"></i>
					{/if}
				</button>
			</form>
		{/if}
		
		<!-- Step 2: OTP Verification -->
		{#if step === 2}
			<form class="form" onsubmit={submitStep2}>
				<p class="otp-info">
					Kami telah mengirim kode verifikasi ke<br/>
					<strong>{email}</strong>
				</p>
				
				<div class="form-group">
					<label for="otp">Kode Verifikasi</label>
					<input
						type="text"
						id="otp"
						bind:value={otp}
						placeholder="123456"
						maxlength="6"
						class="otp-input"
						required
					/>
				</div>
				
				<div class="resend-row">
					{#if resendCooldown > 0}
						<span class="resend-cooldown">Kirim ulang dalam {resendCooldown}s</span>
					{:else}
						<button type="button" class="btn-link" onclick={resendOTP} disabled={loading}>
							Kirim ulang kode
						</button>
					{/if}
				</div>
				
				<button type="submit" class="btn btn-primary btn-full" disabled={loading}>
					{#if loading}
						<i class="bi bi-arrow-repeat spin"></i>
						Memverifikasi...
					{:else}
						Verifikasi
						<i class="bi bi-check-lg"></i>
					{/if}
				</button>
				
				<button type="button" class="btn btn-ghost btn-full" onclick={() => step = 1}>
					<i class="bi bi-arrow-left"></i>
					Kembali
				</button>
			</form>
		{/if}
		
		<!-- Step 3: Username & Slug -->
		{#if step === 3}
			<form class="form" onsubmit={submitStep3}>
				<div class="form-group">
					<label for="username">Username</label>
					<input
						type="text"
						id="username"
						bind:value={username}
						oninput={generateSlug}
						placeholder="johndoe"
						required
					/>
				</div>
				
				<div class="form-group">
					<label for="slug">URL Slug</label>
					<div class="slug-input">
						<span class="slug-prefix">linkmy.deepkernel.site/</span>
						<input
							type="text"
							id="slug"
							bind:value={slug}
							placeholder={username || 'username'}
							pattern="[a-z0-9-]+"
						/>
					</div>
					<span class="input-hint">Boleh kosong untuk menggunakan username</span>
				</div>
				
				<button type="submit" class="btn btn-primary btn-full" disabled={loading}>
					{#if loading}
						<i class="bi bi-arrow-repeat spin"></i>
						Membuat akun...
					{:else}
						<i class="bi bi-rocket-takeoff"></i>
						Buat Akun
					{/if}
				</button>
			</form>
		{/if}
		
		<p class="login-link">
			Sudah punya akun? <a href="/login">Masuk</a>
		</p>
	</div>
</div>

<style>
	.register-page {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: var(--space-xl);
		background: var(--color-bg);
	}
	
	.register-container {
		width: 100%;
		max-width: 420px;
	}
	
	.register-header {
		text-align: center;
		margin-bottom: var(--space-2xl);
	}
	
	.logo {
		display: inline-flex;
		align-items: center;
		gap: var(--space-sm);
		font-size: 1.5rem;
		font-weight: 700;
		color: var(--color-text);
		margin-bottom: var(--space-lg);
	}
	
	.logo i {
		color: var(--color-primary);
		font-size: 1.75rem;
	}
	
	.register-header h1 {
		font-size: 1.75rem;
		margin-bottom: var(--space-xs);
	}
	
	.register-header p {
		color: var(--color-text-muted);
	}
	
	/* Steps */
	.steps {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: var(--space-sm);
		margin-bottom: var(--space-2xl);
	}
	
	.step {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: var(--space-xs);
	}
	
	.step-num {
		width: 36px;
		height: 36px;
		display: flex;
		align-items: center;
		justify-content: center;
		border: 2px solid var(--color-border);
		border-radius: 50%;
		font-weight: 600;
		color: var(--color-text-muted);
		transition: all var(--transition-fast);
	}
	
	.step.active .step-num {
		border-color: var(--color-primary);
		color: var(--color-primary);
	}
	
	.step.completed .step-num {
		background: var(--color-primary);
		border-color: var(--color-primary);
		color: white;
	}
	
	.step-label {
		font-size: 0.75rem;
		color: var(--color-text-muted);
	}
	
	.step.active .step-label {
		color: var(--color-text);
	}
	
	.step-line {
		width: 40px;
		height: 2px;
		background: var(--color-border);
		margin-bottom: 20px;
	}
	
	.step-line.active {
		background: var(--color-primary);
	}
	
	/* Form */
	.form {
		background: var(--color-bg-secondary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-xl);
		padding: var(--space-2xl);
		margin-bottom: var(--space-xl);
	}
	
	.form-group {
		margin-bottom: var(--space-lg);
	}
	
	.form-group label {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		margin-bottom: var(--space-sm);
		color: var(--color-text-secondary);
	}
	
	.form-group input {
		width: 100%;
		padding: var(--space-md);
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		color: var(--color-text);
		font-size: 1rem;
	}
	
	.form-group input:focus {
		outline: none;
		border-color: var(--color-primary);
	}
	
	.otp-info {
		text-align: center;
		margin-bottom: var(--space-xl);
		color: var(--color-text-secondary);
		line-height: 1.6;
	}
	
	.otp-input {
		text-align: center;
		font-size: 1.5rem !important;
		letter-spacing: 0.5em;
		font-weight: 600;
	}
	
	.resend-row {
		text-align: center;
		margin-bottom: var(--space-lg);
	}
	
	.resend-cooldown {
		color: var(--color-text-muted);
		font-size: 0.875rem;
	}
	
	.btn-link {
		background: none;
		border: none;
		color: var(--color-primary);
		cursor: pointer;
		font-size: 0.875rem;
	}
	
	.btn-link:hover {
		text-decoration: underline;
	}
	
	.slug-input {
		display: flex;
		align-items: center;
		background: var(--color-bg-tertiary);
		border: 1px solid var(--color-border);
		border-radius: var(--radius-lg);
		overflow: hidden;
	}
	
	.slug-prefix {
		padding: var(--space-md);
		background: var(--color-bg-secondary);
		color: var(--color-text-muted);
		font-size: 0.75rem;
		white-space: nowrap;
	}
	
	.slug-input input {
		border: none;
		border-radius: 0;
	}
	
	.input-hint {
		display: block;
		font-size: 0.75rem;
		color: var(--color-text-muted);
		margin-top: var(--space-xs);
	}
	
	.btn-full {
		width: 100%;
		justify-content: center;
		margin-bottom: var(--space-sm);
	}
	
	/* Alert */
	.alert {
		padding: var(--space-md);
		border-radius: var(--radius-md);
		margin-bottom: var(--space-lg);
		display: flex;
		align-items: center;
		gap: var(--space-sm);
	}
	
	.alert-error {
		background: rgba(239, 68, 68, 0.1);
		border: 1px solid rgba(239, 68, 68, 0.3);
		color: #ef4444;
	}
	
	.login-link {
		text-align: center;
		color: var(--color-text-secondary);
	}
	
	.login-link a {
		color: var(--color-primary);
		font-weight: 500;
	}
	
	.spin {
		animation: spin 1s linear infinite;
	}
	
	@keyframes spin {
		to { transform: rotate(360deg); }
	}
</style>
