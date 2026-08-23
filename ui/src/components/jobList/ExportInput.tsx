import { useFieldArray, useFormContext } from 'react-hook-form';
import type { JobSpecFormValues } from '../../schemas/jobSpec';
import { COMMON_LABELS } from '../../constants/common';
import { CREATE_JOB_TEXTS } from '../../constants/createJob';
import AppButton from '../common/AppButton';
import AppInput from '../common/AppInput';

/**
 * The "Exports" section of the create-job form: a required, repeatable
 * list of output file names. The storage backend is fixed to S3
 * internally - the user only ever names the file.
 */
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
            <AppInput
              placeholder={CREATE_JOB_TEXTS.EXPORT_PATH_PLACEHOLDER}
              aria-label={CREATE_JOB_TEXTS.EXPORT_PATH_ARIA}
              {...register(`exports.${i}.path`)}
              className="flex-1"
            />
            <AppButton
              type="button"
              variant="ghost"
              onClick={() => remove(i)}
              aria-label={CREATE_JOB_TEXTS.REMOVE_EXPORT_ARIA}
              className="text-neutral-400! hover:text-danger-600!"
            >
              {COMMON_LABELS.REMOVE}
            </AppButton>
          </div>
          {errors.exports?.[i]?.path && (
            <p className="text-xs text-danger-600 mt-1">{errors.exports[i]?.path?.message}</p>
          )}
        </div>
      ))}
      {errors.exports?.message && <p className="text-xs text-danger-600">{errors.exports.message}</p>}
      <AppButton
        type="button"
        variant="ghost"
        onClick={() => append({ type: 's3', path: '' })}
        className="text-brand-600! hover:text-brand-700! hover:bg-transparent!"
      >
        {CREATE_JOB_TEXTS.ADD_EXPORT}
      </AppButton>
    </div>
  );
}
