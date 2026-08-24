import AppButton from './AppButton';

import { ConfirmModalProps } from '../../types/common';
import { COMMON_LABELS } from '../../constants/common';

/**
 * Generic yes/no confirmation dialog, reused for anything that needs a
 * "are you sure?" step before acting (deleting or cancelling a job).
 * Renders nothing while closed.
 */
export default function ConfirmModal({ isOpen, title, message, onConfirm, onCancel }: ConfirmModalProps) {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-[100]">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-md p-6">
        <h3 className="text-lg font-semibold text-neutral-900 mb-2">{title}</h3>
        <p className="text-sm text-neutral-500 mb-6">{message}</p>
        <div className="flex justify-end gap-3">
          <AppButton variant="secondary" onClick={onCancel}>{COMMON_LABELS.CANCEL}</AppButton>
          <AppButton variant="danger" onClick={onConfirm}>{COMMON_LABELS.CONFIRM}</AppButton>
        </div>
      </div>
    </div>
  );
}
