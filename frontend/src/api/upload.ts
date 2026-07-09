import http from './http';

interface UploadResponse {
    url: string;
}

export async function uploadFile(file: File): Promise<UploadResponse> {
    const formData = new FormData();
    formData.append('file', file);
    const { data } = await http.post<UploadResponse>('/upload', formData, {
        headers: {
            'Content-Type': 'multipart/form-data',
        },
    });
    return data;
}