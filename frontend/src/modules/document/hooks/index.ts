import { AxiosError } from 'axios';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { DocumentService } from '../services';
import type { ApproveDocumentRequest, Document, RejectDocumentRequest } from '../types';
import { toast } from 'react-hot-toast';
import { useNavigate } from '@tanstack/react-router';

export const useDocuments = () => {
  return useQuery({
    queryKey: ['documents'],
    queryFn: DocumentService.getDocuments,
  });
};

export const useDocument = (id: string) => {
  return useQuery({
    queryKey: ['documents', id],
    queryFn: () => DocumentService.getDocument(id),
    enabled: !!id,
  });
};

export const useUploadDocument = () => {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: DocumentService.uploadDocument,
    onSuccess: (res) => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      toast.success('Document uploaded successfully');
      navigate({ to: `/documents/${res.id}` });
    },
    onError: (error: AxiosError<unknown>) => {
      const message = (error.response?.data as { message?: string } | undefined)?.message;
      toast.error(message || 'Failed to upload document');
    },
  });
};

export const useUpdateDocument = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: Document }) =>
      DocumentService.updateDocument(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      if (variables?.id) {
        queryClient.invalidateQueries({ queryKey: ['documents', variables.id] });
      }
      toast.success('Document saved successfully');
    },
    onError: (error: AxiosError<unknown>) => {
      const message = (error.response?.data as { message?: string } | undefined)?.message;
      toast.error(message || 'Failed to save document');
    },
  });
};

export const useApproveDocument = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data?: ApproveDocumentRequest }) =>
      DocumentService.approveDocument(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      if (variables?.id) {
        queryClient.invalidateQueries({ queryKey: ['documents', variables.id] });
      }
      toast.success('Document approved');
    },
    onError: (error: AxiosError<unknown>) => {
      const message = (error.response?.data as { message?: string } | undefined)?.message;
      toast.error(message || 'Failed to approve document');
    },
  });
};

export const useRejectDocument = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: RejectDocumentRequest }) =>
      DocumentService.rejectDocument(id, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['documents'] });
      if (variables?.id) {
        queryClient.invalidateQueries({ queryKey: ['documents', variables.id] });
      }
      toast.success('Document rejected');
    },
    onError: (error: AxiosError<unknown>) => {
      const message = (error.response?.data as { message?: string } | undefined)?.message;
      toast.error(message || 'Failed to reject document');
    },
  });
};