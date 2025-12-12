// API Configuration and Client
const API_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:3000';

interface ApiResponse<T> {
	success?: boolean;
	data?: T;
	error?: string;
	message?: string;
}

interface User {
	id: number;
	username: string;
	email: string;
	is_verified: boolean;
	created_at: string;
}

interface AuthResponse {
	user: User;
	access_token: string;
	refresh_token: string;
	expires_in: number;
}

interface Profile {
	id: number;
	user_id: number;
	slug: string;
	name: string;
	title?: string;
	bio?: string;
	avatar: string;
	is_active: boolean;
	display_order: number;
	link_count?: number;
	total_clicks?: number;
}

interface Link {
	id: number;
	profile_id: number;
	category_id?: number;
	title: string;
	url: string;
	icon: string;
	position: number;
	clicks: number;
	is_active: boolean;
}

interface Category {
	id: number;
	profile_id: number;
	name: string;
	icon: string;
	color: string;
	position: number;
	is_expanded: boolean;
}

interface Theme {
	id: number;
	profile_id: number;
	bg_type: string;
	bg_value?: string;
	button_style: string;
	button_color: string;
	text_color: string;
	font: string;
	layout: string;
	container_style: string;
	enable_animations: boolean;
	enable_glass_effect: boolean;
	shadow_intensity: string;
	boxed_enabled: boolean;
	boxed_outer_bg_type?: string;
	boxed_outer_bg_value?: string;
	boxed_container_bg: string;
	boxed_max_width: number;
	boxed_radius: number;
	boxed_shadow: boolean;
}

interface PublicProfile {
	profile: Profile;
	theme: Theme;
	categories: Category[];
	links: Link[];
	is_verified: boolean;
}

// Token storage
function getToken(): string | null {
	if (typeof window === 'undefined') return null;
	return localStorage.getItem('access_token');
}

function setTokens(access: string, refresh: string): void {
	localStorage.setItem('access_token', access);
	localStorage.setItem('refresh_token', refresh);
}

function clearTokens(): void {
	localStorage.removeItem('access_token');
	localStorage.removeItem('refresh_token');
}

function getRefreshToken(): string | null {
	if (typeof window === 'undefined') return null;
	return localStorage.getItem('refresh_token');
}

// API client
async function request<T>(
	endpoint: string,
	options: RequestInit = {}
): Promise<ApiResponse<T>> {
	const token = getToken();
	
	const headers: HeadersInit = {
		'Content-Type': 'application/json',
		...options.headers
	};
	
	if (token) {
		(headers as Record<string, string>)['Authorization'] = `Bearer ${token}`;
	}
	
	try {
		const response = await fetch(`${API_URL}${endpoint}`, {
			...options,
			headers
		});
		
		const data = await response.json();
		
		if (!response.ok) {
			return { error: data.message || 'Request failed' };
		}
		
		return data;
	} catch (err) {
		return { error: 'Network error' };
	}
}

// Auth API
export const auth = {
	async register(username: string, email: string, password: string): Promise<ApiResponse<AuthResponse>> {
		const res = await request<AuthResponse>('/api/v1/auth/register', {
			method: 'POST',
			body: JSON.stringify({ username, email, password })
		});
		if (res.data) {
			setTokens(res.data.access_token, res.data.refresh_token);
		}
		return res;
	},
	
	async login(email: string, password: string): Promise<ApiResponse<AuthResponse>> {
		const res = await request<AuthResponse>('/api/v1/auth/login', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
		if (res.data) {
			setTokens(res.data.access_token, res.data.refresh_token);
		}
		return res;
	},
	
	async logout(): Promise<void> {
		const refreshToken = getRefreshToken();
		await request('/api/v1/auth/logout', {
			method: 'POST',
			body: JSON.stringify({ refresh_token: refreshToken })
		});
		clearTokens();
	},
	
	async getCurrentUser(): Promise<ApiResponse<User>> {
		return request<User>('/api/v1/me');
	},
	
	isAuthenticated(): boolean {
		return !!getToken();
	},
	
	// OTP Registration Flow
	async sendOTP(email: string, password: string): Promise<ApiResponse<{ message: string; email: string }>> {
		return request<{ message: string; email: string }>('/api/v1/auth/send-otp', {
			method: 'POST',
			body: JSON.stringify({ email, password })
		});
	},
	
	async verifyOTP(email: string, otp: string): Promise<ApiResponse<{ message: string; verified: boolean }>> {
		return request<{ message: string; verified: boolean }>('/api/v1/auth/verify-otp', {
			method: 'POST',
			body: JSON.stringify({ email, otp })
		});
	},
	
	async completeRegistration(email: string, password: string, otp: string, username: string, slug: string): Promise<ApiResponse<AuthResponse>> {
		const res = await request<AuthResponse>('/api/v1/auth/complete-registration', {
			method: 'POST',
			body: JSON.stringify({ email, password, otp, username, slug })
		});
		if (res.data) {
			setTokens(res.data.access_token, res.data.refresh_token);
		}
		return res;
	}
};

// Profile API
export const profiles = {
	async getPublic(slug: string): Promise<ApiResponse<PublicProfile>> {
		return request<PublicProfile>(`/api/v1/p/${slug}`);
	},
	
	async getAll(): Promise<ApiResponse<Profile[]>> {
		return request<Profile[]>('/api/v1/profiles');
	},
	
	async get(id: number): Promise<ApiResponse<Profile>> {
		return request<Profile>(`/api/v1/profiles/${id}`);
	},
	
	async create(data: Partial<Profile>): Promise<ApiResponse<Profile>> {
		return request<Profile>('/api/v1/profiles', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	
	async update(id: number, data: Partial<Profile>): Promise<ApiResponse<Profile>> {
		return request<Profile>(`/api/v1/profiles/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},
	
	async delete(id: number): Promise<ApiResponse<void>> {
		return request<void>(`/api/v1/profiles/${id}`, { method: 'DELETE' });
	}
};

// Links API
export const links = {
	async getByProfile(profileId: number): Promise<ApiResponse<Link[]>> {
		return request<Link[]>(`/api/v1/profiles/${profileId}/links`);
	},
	
	async create(profileId: number, data: Partial<Link>): Promise<ApiResponse<Link>> {
		return request<Link>(`/api/v1/profiles/${profileId}/links`, {
			method: 'POST',
			body: JSON.stringify(data)
		});
	},
	
	async update(id: number, data: Partial<Link>): Promise<ApiResponse<Link>> {
		return request<Link>(`/api/v1/links/${id}`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	},
	
	async delete(id: number): Promise<ApiResponse<void>> {
		return request<void>(`/api/v1/links/${id}`, { method: 'DELETE' });
	},
	
	async trackClick(id: number): Promise<ApiResponse<{ url: string }>> {
		return request<{ url: string }>(`/api/v1/click/${id}`, { method: 'POST' });
	}
};

// Theme API
export const themes = {
	async get(profileId: number): Promise<ApiResponse<Theme>> {
		return request<Theme>(`/api/v1/profiles/${profileId}/theme`);
	},
	
	async update(profileId: number, data: Partial<Theme>): Promise<ApiResponse<Theme>> {
		return request<Theme>(`/api/v1/profiles/${profileId}/theme`, {
			method: 'PUT',
			body: JSON.stringify(data)
		});
	}
};

export type { User, Profile, Link, Category, Theme, PublicProfile, AuthResponse, ApiResponse };
