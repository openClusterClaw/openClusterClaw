import apiClient from './client';

// API Response wrapper
interface ApiResponse<T> {
	code: number;
	message: string;
	data: T;
}

// Generate OTP secret response
export interface GenerateOTPResponse {
	secret: string;
	qr_code: string;
}

// Enable OTP request
export interface EnableOTPRequest {
	code: string;
}

// Enable OTP response
export interface EnableOTPResponse {
	backup_codes: string[];
}

// Disable OTP request
export interface DisableOTPRequest {
	code: string;
}

// Verify OTP request
export interface VerifyOTPRequest {
	temp_token: string;
	code: string;
}

// OTP status
export interface OTPStatus {
	otp_enabled: boolean;
}

// OTP API
export const otpApi = {
	// Generate OTP secret
	generateSecret: async (): Promise<GenerateOTPResponse> => {
		const response = await apiClient.post<ApiResponse<GenerateOTPResponse>>('/auth/otp/generate');
		return response.data.data;
	},

	// Enable OTP
	enableOTP: async (req: EnableOTPRequest): Promise<EnableOTPResponse> => {
		const response = await apiClient.post<ApiResponse<EnableOTPResponse>>('/auth/otp/enable', req);
		return response.data.data;
	},

	// Disable OTP
	disableOTP: async (req: DisableOTPRequest): Promise<void> => {
		const response = await apiClient.post<ApiResponse<void>>('/auth/otp/disable', req);
		return response.data.data;
	},

	// Verify OTP
	verifyOTP: async (req: VerifyOTPRequest) => {
		const response = await apiClient.post<ApiResponse<any>>('/auth/otp/verify', req);
		return response.data.data;
	},

	// Get backup codes
	getBackupCodes: async (): Promise<string[]> => {
		const response = await apiClient.get<ApiResponse<{ backup_codes: string[] }>>('/auth/otp/backup');
		return response.data.data.backup_codes;
	},

	// Get OTP status
	getOTPStatus: async (): Promise<OTPStatus> => {
		const response = await apiClient.get<ApiResponse<OTPStatus>>('/auth/otp/status');
		return response.data.data;
	},
};

export default otpApi;
