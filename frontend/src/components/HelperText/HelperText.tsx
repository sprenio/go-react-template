import React, { useState } from 'react';
import { HelpCircle } from 'lucide-react';

interface HelperTextProps {
  helperText?: string;
  className?: string;
  helperIconClassName?: string;
  tooltipClassName?: string;
}

export const HelperText: React.FC<HelperTextProps> = ({
  helperText,
  className = '',
  helperIconClassName = '',
  tooltipClassName = '',
}) => {
  const [showHelper, setShowHelper] = useState(false);

  return (
    <span
      className={`helper-wrapper ${className}`}
      style={{ display: 'inline-flex', alignItems: 'center', gap: 4, position: 'relative' }}
    >
      {/* Ikonka helpera */}
      {helperText && (
        <span
          onMouseEnter={() => setShowHelper(true)}
          onMouseLeave={() => setShowHelper(false)}
          onClick={(e: React.MouseEvent) => {
            e.preventDefault();
            setShowHelper((prev) => !prev);
          }}
          className={`helper-icon ${helperIconClassName}`}
          style={{ cursor: 'pointer', display: 'inline-flex', alignItems: 'center' }}
        >
          <HelpCircle size={16} />
        </span>
      )}

      {/* Tooltip */}
      {helperText && showHelper && (
        <div
          className={`helper-tooltip ${tooltipClassName}`}
          style={{
            position: 'absolute',
            top: '100%',
            left: 0,
            background: 'rgba(0,0,0,0.85)',
            color: 'white',
            padding: '8px 12px',
            borderRadius: 6,
            fontSize: '0.875rem',
            whiteSpace: 'normal',
            width: 'max-content',
            maxWidth: '320px',
            zIndex: 1000,
            marginTop: 8,
            boxShadow: '0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06)',
            lineHeight: '1.4',
          }}
          dangerouslySetInnerHTML={{ __html: helperText }}
        ></div>
      )}
    </span>
  );
};
