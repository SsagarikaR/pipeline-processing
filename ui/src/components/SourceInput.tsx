import { useFieldArray, useFormContext } from 'react-hook-form';
import type { JobSpecFormValues } from '../schemas/jobSpec';

export default function SourceInput() {
  const {
    register,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useFormContext<JobSpecFormValues>();
  const { fields, append, remove } = useFieldArray({ control, name: 'sources' });
  const types = watch('sources');

  function handleUpload(index: number, file: File) {
    const reader = new FileReader();
    reader.onload = (event) => {
      if (typeof event.target?.result === 'string') {
        setValue(`sources.${index}.path`, event.target.result, {
          shouldValidate: true,
          shouldDirty: true,
        });
      }
    };
    reader.readAsDataURL(file);
  }

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">
        Sources <span className="text-danger-600">*</span>
      </label>
      {fields.map((field, i) => (
        <div key={field.id}>
          <div className="flex gap-2">
            <select
              {...register(`sources.${i}.type`)}
              className="border border-neutral-300 rounded-lg px-3 py-2 text-sm"
            >
              <option value="csv">CSV</option>
              <option value="json">JSON</option>
            </select>
            <div className="flex-1 flex gap-2">
              <input
                type="text"
                placeholder="/path/to/file or Data URI"
                {...register(`sources.${i}.path`)}
                className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
              />
              <label className="cursor-pointer border border-neutral-300 rounded-lg px-3 py-2 text-sm bg-neutral-50 hover:bg-neutral-100 flex items-center shrink-0">
                <span>Upload</span>
                <input
                  type="file"
                  className="hidden"
                  accept={types?.[i]?.type === 'json' ? '.json' : '.csv'}
                  onChange={(e) => {
                    const file = e.target.files?.[0];
                    if (file) handleUpload(i, file);
                    e.target.value = ''; // reset
                  }}
                />
              </label>
            </div>
            <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">
              ✕
            </button>
          </div>
          {errors.sources?.[i]?.path && (
            <p className="text-xs text-danger-600 mt-1">{errors.sources[i]?.path?.message}</p>
          )}
        </div>
      ))}
      {errors.sources?.message && <p className="text-xs text-danger-600">{errors.sources.message}</p>}
      <button
        type="button"
        onClick={() => append({ type: 'csv', path: '' })}
        className="text-sm text-brand-600 hover:text-brand-700 font-medium"
      >
        + Add source
      </button>
    </div>
  );
}
