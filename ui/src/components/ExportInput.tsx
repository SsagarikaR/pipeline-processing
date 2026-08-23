import { useFieldArray, useFormContext } from 'react-hook-form';
import type { JobSpecFormValues } from '../schemas/jobSpec';

export default function ExportInput() {
  const {
    register,
    control,
    formState: { errors },
  } = useFormContext<JobSpecFormValues>();
  const { fields, append, remove } = useFieldArray({ control, name: 'exports' });

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">
        Exports <span className="text-danger-600">*</span>
      </label>
      {fields.map((field, i) => (
        <div key={field.id}>
          <div className="flex gap-2 items-center">
            <input
              placeholder="e.g., results.json"
              {...register(`exports.${i}.path`)}
              className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
            />
            <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">
              ✕
            </button>
          </div>
          {errors.exports?.[i]?.path && (
            <p className="text-xs text-danger-600 mt-1">{errors.exports[i]?.path?.message}</p>
          )}
        </div>
      ))}
      {errors.exports?.message && <p className="text-xs text-danger-600">{errors.exports.message}</p>}
      <button
        type="button"
        onClick={() => append({ type: 's3', path: '' })}
        className="text-sm text-brand-600 hover:text-brand-700 font-medium"
      >
        + Add export
      </button>
    </div>
  );
}
