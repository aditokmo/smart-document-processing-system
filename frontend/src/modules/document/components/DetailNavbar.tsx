import { useEffect, useRef, useState } from 'react';
import { ArrowLeft } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface DetailNavbarProps {
  filename: string;
  isEditing: boolean;
  onEdit: () => void;
  onCancel: () => void;
  onSave: () => void;
  isSaving: boolean;
  saveDisabled: boolean;
}

export function DetailNavbar({ filename, isEditing, onEdit, onCancel, onSave, isSaving, saveDisabled }: DetailNavbarProps) {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);

  const handleBack = () => {
    window.history.back();
  };

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (isOpen && containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [isOpen]);

  return (
    <div className="bg-white border-b border-gray-200 sticky top-0 z-10">
      <div className="px-6 py-4 flex items-center justify-between gap-4 max-w-6xl mx-auto relative" ref={containerRef}>
        <div className="flex items-center gap-4">
          <Button
            onClick={handleBack}
            variant="ghost"
            size="icon"
            className="text-gray-600 hover:text-gray-900"
          >
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <h1 className="text-xl font-semibold text-gray-900">{filename}</h1>
        </div>

        <div className="flex items-center gap-2">
          {isEditing ? (
            <>
              <Button onClick={onCancel} variant="outline" disabled={isSaving}>
                Cancel
              </Button>
              <Button onClick={onSave} disabled={saveDisabled}>
                {isSaving ? 'Saving...' : 'Save'}
              </Button>
            </>
          ) : (
            <Button onClick={onEdit}>Edit</Button>
          )}
        </div>
      </div>
    </div>
  );
}
