import { useParams } from '@tanstack/react-router';
import { useDocument, useApproveDocument, useRejectDocument, useUpdateDocument } from '../hooks';
import { useForm, useWatch } from 'react-hook-form';
import { useState, useEffect } from 'react';
import type { Document as DocType, LineItem as LineItemType } from '../types';

type EditableDoc = DocType & { tax_rate?: number };
import {
  DetailNavbar,
  DocumentHeader,
  DocumentDetails,
  FinancialSummary,
  LineItemsTable,
  ValidationIssues,
  ActionButtons,
} from '../components';

export default function DocumentDetail() {
  const { id } = useParams({ from: '/document/$id' });
  const { data: doc, isLoading } = useDocument(id);
  const approveMutation = useApproveDocument();
  const rejectMutation = useRejectDocument();
  const updateMutation = useUpdateDocument();
  
  const {
    setValue,
    reset,
    control,
    formState: { isDirty },
  } = useForm<EditableDoc>();

  const [isEditing, setIsEditing] = useState(false);
  const [rejectReason, setRejectReason] = useState('');
  const watchedFormData = useWatch({ control });
  const formData = watchedFormData as EditableDoc;

  useEffect(() => {
    if (doc) {
      reset({
        ...doc,
        tax_rate: doc.subtotal ? (doc.tax / doc.subtotal) * 100 : 0,
      } as EditableDoc);
    }
  }, [doc, reset]);

  if (isLoading) return <div className="flex justify-center items-center h-64">Loading...</div>;
  if (!doc) return <div className="text-center text-red-500">Document not found</div>;

  const handleApprove = () => {
    const data = { corrections: formData };
    approveMutation.mutate({ id, data });
  };

  const handleReject = () => {
    if (!rejectReason?.trim()) return;
    rejectMutation.mutate({ id, data: { reason: rejectReason } });
  };

  const handleSave = () => {
    if (!doc) return;

    const updatedDocument = {
      ...doc,
      ...formData,
      line_items: lineItems,
      subtotal: formData.subtotal ?? doc.subtotal,
      tax: formData.tax ?? doc.tax,
      total: formData.total ?? doc.total,
    } as DocType;

    updateMutation.mutate(
      { id, data: updatedDocument },
      {
        onSuccess: () => {
          setIsEditing(false);
        },
      }
    );
  };

  const handleCancel = () => {
    if (!doc) return;
    reset({
      ...doc,
      tax_rate: doc.subtotal ? (doc.tax / doc.subtotal) * 100 : 0,
    } as EditableDoc);
    setRejectReason('');
    setIsEditing(false);
  };

  const isEditable = isEditing;
  const lineItems = formData.line_items ?? doc.line_items;

  const calculateTotals = (next: Partial<EditableDoc>) => {
    const subtotal = next.line_items?.reduce((sum, item) => sum + (item.total || 0), 0) ?? (next.subtotal ?? doc.subtotal);
    const taxRate = next.tax_rate ?? (doc.subtotal ? (doc.tax / doc.subtotal) * 100 : 0);
    const tax = subtotal * (taxRate / 100);
    const total = next.total ?? subtotal + tax;
    return { ...next, subtotal, tax, total };
  };

  const updateFormData = (key: keyof EditableDoc, value: unknown) => {
    const next = { ...formData, [key]: value } as Partial<EditableDoc>;
    if ((key === 'subtotal' || key === 'tax_rate') && formData.total !== undefined) {
      delete next.total;
    }
    const updated = calculateTotals(next);
    Object.entries(updated).forEach(([k, v]) => {
      setValue(k as keyof EditableDoc, v, { shouldDirty: true });
    });
  };

  const updateLineItem = (index: number, key: keyof LineItemType, value: unknown) => {
    const updatedItems = [...lineItems];
    updatedItems[index] = { ...updatedItems[index], [key]: value };

    if (key === 'quantity' || key === 'price') {
      updatedItems[index].total = (updatedItems[index].quantity || 0) * (updatedItems[index].price || 0);
    }

    const next = { ...formData, line_items: updatedItems };
    const updated = calculateTotals(next as Partial<EditableDoc>);
    Object.entries(updated).forEach(([k, v]) => {
      setValue(k as keyof EditableDoc, v, { shouldDirty: true });
    });
  };

  const addLineItem = () => {
    const next = {
      ...formData,
      line_items: [...(formData.line_items || doc.line_items), { description: '', quantity: 0, price: 0, total: 0 } as LineItemType],
    };
    const updated = calculateTotals(next as Partial<EditableDoc>);
    Object.entries(updated).forEach(([k, v]) => {
      setValue(k as keyof EditableDoc, v, { shouldDirty: true });
    });
  };

  return (
    <>
      <DetailNavbar
        filename={formData.filename || formData.document_number || doc.filename || doc.document_number || 'Document'}
        isEditing={isEditing}
        onEdit={() => setIsEditing(true)}
        onCancel={handleCancel}
        onSave={handleSave}
        isSaving={updateMutation.isPending}
        saveDisabled={updateMutation.isPending || !isDirty}
      />
      <div className="p-6 max-w-6xl mx-auto space-y-6">
        <DocumentHeader
          doc={doc}
          documentNumber={formData.document_number || doc.document_number || ''}
          supplierName={formData.supplier_name || doc.supplier_name || ''}
        />

        <div className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <DocumentDetails
              doc={doc}
              formData={formData}
              isEditable={isEditable}
              onUpdate={updateFormData}
            />

            <FinancialSummary
              doc={doc}
              formData={formData}
              isEditable={isEditable}
              onUpdate={updateFormData}
            />
          </div>

          <LineItemsTable
            doc={doc}
            lineItems={lineItems}
            isEditable={isEditable}
            onUpdateLineItem={updateLineItem}
            onAddLineItem={addLineItem}
          />

          <ValidationIssues issues={doc.issues || []} />

          <ActionButtons
            status={doc.status}
            rejectReason={rejectReason}
            onRejectReasonChange={setRejectReason}
            onApprove={handleApprove}
            onReject={handleReject}
            isApprovePending={approveMutation.isPending}
            isRejectPending={rejectMutation.isPending}
            hasChanges={isDirty}
            hasIssues={Boolean(doc.issues && doc.issues.length > 0)}
            isEditing={isEditing}
          />
        </div>
      </div>
    </>
  );
}