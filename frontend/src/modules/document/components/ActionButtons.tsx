import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { CheckCircle, XCircle } from 'lucide-react';

interface ActionButtonsProps {
  status: string;
  rejectReason: string;
  onRejectReasonChange: (reason: string) => void;
  onApprove: () => void;
  onReject: () => void;
  isApprovePending: boolean;
  isRejectPending: boolean;
  hasChanges?: boolean;
  hasIssues: boolean;
  isEditing: boolean;
}

export function ActionButtons({
  status,
  rejectReason,
  onRejectReasonChange,
  onApprove,
  onReject,
  isApprovePending,
  isRejectPending,
  hasChanges,
  hasIssues,
  isEditing,
}: ActionButtonsProps) {
  if (status !== 'needs_review') return null;

  const disableApprove = isApprovePending || isEditing || (hasIssues && !hasChanges);

  return (
    <div className="bg-gray-50 rounded-lg p-6 border border-gray-200 space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 items-start">
        <div className="md:col-span-2">
          <label className="block text-sm font-medium mb-2 text-gray-700">
            Reason for Rejection (if rejecting)
          </label>
          <Textarea
            placeholder="Explain why this document is being rejected..."
            value={rejectReason}
            onChange={(e) => onRejectReasonChange(e.target.value)}
            className="min-h-20 resize-none"
          />
        </div>

        <div className="flex flex-col gap-3 md:pt-7">
          <Button
            onClick={onApprove}
            disabled={disableApprove}
            className="bg-green-600 hover:bg-green-700 text-white w-full"
            size="lg"
          >
            <CheckCircle className="mr-2 h-4 w-4" />
            Approve
          </Button>
          <Button
            onClick={onReject}
            disabled={isRejectPending || !rejectReason.trim()}
            variant="destructive"
            className="w-full"
            size="lg"
          >
            <XCircle className="mr-2 h-4 w-4" />
            Reject
          </Button>
        </div>
      </div>

      {isEditing && (
        <p className="text-sm text-gray-500 md:col-span-3">
          💡 Save or cancel your edits.
        </p>
      )}
      {hasIssues && !hasChanges && !isEditing && (
        <p className="text-sm text-gray-500 md:col-span-3">
          💡 Fix the validation issues to be able to approve document.
        </p>
      )}
      {!hasIssues && !isEditing && (
        <p className="text-sm text-gray-500 md:col-span-3">
          💡 No unresolved validation issues. You can approve this document.
        </p>
      )}
      {!rejectReason.trim() && (
        <p className="text-sm text-gray-500 md:col-span-3">
          💡 Add a rejection reason to enable the reject button
        </p>
      )}
    </div>
  );
}
