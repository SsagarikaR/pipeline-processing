import { useFieldArray, useFormContext } from 'react-hook-form';
import type { JobSpecFormValues } from '../../schemas/jobSpec';
import { COMMON_LABELS } from '../../constants/common';
import { CREATE_JOB_TEXTS } from '../../constants/createJob';
import AppButton from '../common/AppButton';
import AppInput from '../common/AppInput';

/**
 * The "Sources" section of the create-job form: a repeatable list of
 * input files (CSV or JSON), each with a type, a path, and a file
 * upload shortcut that fills the path in for you. Reads and writes
 * directly into the shared CreateJobModal form via useFormContext.
 */
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

  /** Reads a chosen file and fills the row's path with it as a data URI. */
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
        {CREATE_JOB_TEXTS.SOURCES_LABEL} <span className="text-danger-600">*</span>
      </label>
      {fields.map((field, i) => (
        <div key={field.id}>
          <div className="flex gap-2">
            <select
              {...register(`sources.${i}.type`)}
              aria-label={CREATE_JOB_TEXTS.SOURCE_TYPE_ARIA}
              className="border border-neutral-300 rounded-lg px-3 py-2 text-sm"
            >
              <option value="csv">CSV</option>
              <option value="json">JSON</option>
            </select>
            <div className="flex-1 flex gap-2">
              <AppInput
                type="text"
                placeholder={CREATE_JOB_TEXTS.SOURCE_PATH_PLACEHOLDER}
                aria-label={CREATE_JOB_TEXTS.SOURCE_PATH_ARIA}
                {...register(`sources.${i}.path`)}
                className="flex-1"
              />
              <label className="cursor-pointer border border-neutral-300 rounded-lg px-3 py-2 text-sm bg-neutral-50 hover:bg-neutral-100 flex items-center shrink-0">
                <span>{CREATE_JOB_TEXTS.SOURCE_UPLOAD}</span>
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
            <AppButton
              type="button"
              variant="ghost"
              onClick={() => remove(i)}
              aria-label={CREATE_JOB_TEXTS.REMOVE_SOURCE_ARIA}
              className="text-neutral-400! hover:text-danger-600!"
            >
              {COMMON_LABELS.REMOVE}
            </AppButton>
          </div>
          {errors.sources?.[i]?.path && (
            <p className="text-xs text-danger-600 mt-1">{errors.sources[i]?.path?.message}</p>
          )}
        </div>
      ))}
      {errors.sources?.message && <p className="text-xs text-danger-600">{errors.sources.message}</p>}
      <AppButton
        type="button"
        variant="ghost"
        onClick={() => append({ type: 'csv', path: '' })}
        className="text-brand-600! hover:text-brand-700! hover:bg-transparent!"
      >
        {CREATE_JOB_TEXTS.ADD_SOURCE}
      </AppButton>
    </div>
  );
}
