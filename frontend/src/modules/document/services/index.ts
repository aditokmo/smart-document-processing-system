import { ApiService } from '../../../api/api';
import type { ApproveDocumentRequest, Document, UploadDocumentRequest, RejectDocumentRequest } from '../types';

export class DocumentService {
  static async uploadDocument(data: UploadDocumentRequest): Promise<Document> {
    const formData = new FormData();
    formData.append('file', data.file);
    return ApiService.post('/documents', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  }

  static async getDocuments(): Promise<Document[]> {
    return ApiService.get('/documents');
  }

  static async getDocument(id: string): Promise<Document> {
    return ApiService.get(`/documents/${id}`);
  }

  static async updateDocument(id: string, data: Document): Promise<Document> {
    return ApiService.put<Document, Document>(`/documents/${id}`, data);
  }

  static async approveDocument(id: string, data?: ApproveDocumentRequest): Promise<void> {
    return ApiService.post(`/documents/${id}/approve`, data);
  }

  static async rejectDocument(id: string, data: RejectDocumentRequest): Promise<void> {
    return ApiService.post(`/documents/${id}/reject`, data);
  }
}