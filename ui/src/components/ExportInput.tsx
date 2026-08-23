import { useFieldArray, useFormContext } from 'react-hook-form';
import type { JobSpecFormValues } from '../schemas/jobSpec';
import { COMMON_LABELS } from '../constants/common';
import { CREATE_JOB_TEXTS } from '../constants/createJob';

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
        {CREATE_JOB_TEXTS.EXPORTS_LABEL} <span className="text-danger-600">*</span>
      </label>
      {fields.map((field, i) => (
        <div key={field.id}>
          <div className="flex gap-2 items-center">
            <input
              placeholder={CREATE_JOB_TEXTS.EXPORT_PATH_PLACEHOLDER}
              {...register(`exports.${i}.path`)}
              className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
            />
            <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">
              {COMMON_LABELS.REMOVE}
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
        {CREATE_JOB_TEXTS.ADD_EXPORT}
      </button>
    </div>
  );
}
