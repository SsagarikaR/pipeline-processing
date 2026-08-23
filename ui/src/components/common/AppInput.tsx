import { forwardRef, useId } from 'react';
import type { AppInputProps } from '../../types/common';

/**
 * Standard text/number input used everywhere in the app instead of a
 * bare `<input>`, so every field shares the same border/padding/text
 * size. Pass `label` to render a properly-associated `<label>` above
 * the input; pass `aria-label` instead when a visible label would be
 * redundant (e.g. a compact repeated-row field whose section already
 * has its own label). Forwards its ref through to the real DOM node,
 * which react-hook-form's `register()` relies on to work correctly.
 */
const AppInput = forwardRef<HTMLInputElement, AppInputProps>(function AppInput(
  { className = '', label, id, ...props },
  ref
) {
  const generatedId = useId();
  const inputId = id ?? generatedId;

  const input = (
    <input
      ref={ref}
      id={inputId}
      className={`border border-neutral-300 rounded-lg px-3 py-2 text-sm ${className}`}
      {...props}
    />
  );

  if (!label) return input;

  return (
    <div>
      <label htmlFor={inputId} className="block text-xs text-neutral-500 mb-1">
        {label}
      </label>
      {input}
    </div>
  );
});

export default AppInput;
