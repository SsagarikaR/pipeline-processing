import { useFieldArray, useFormContext } from 'react-hook-form';
import type { JobSpecFormValues } from '../../schemas/jobSpec';
import { COMMON_LABELS } from '../../constants/common';
import { CREATE_JOB_TEXTS } from '../../constants/createJob';
import AppButton from '../common/AppButton';
import AppInput from '../common/AppInput';

/**
 * The "Transforms" section of the create-job form: an optional,
 * repeatable list of field transforms (uppercase/lowercase) to apply
 * before aggregation.
 */
export default function TransformInput() {
  const {
    register,
    control,
    formState: { errors },
  } = useFormContext<JobSpecFormValues>();
  const { fields, append, remove } = useFieldArray({ control, name: 'transforms' });

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">{CREATE_JOB_TEXTS.TRANSFORMS_LABEL}</label>
      {fields.map((field, i) => (
        <div key={field.id}>
          <div className="flex gap-2 items-center">
            <select
              {...register(`transforms.${i}.name`)}
              aria-label={CREATE_JOB_TEXTS.TRANSFORM_NAME_ARIA}
              className="border border-neutral-300 rounded-lg px-3 py-2 text-sm"
            >
              <option value="uppercase">uppercase</option>
              <option value="lowercase">lowercase</option>
            </select>

            <AppInput
              placeholder={CREATE_JOB_TEXTS.TRANSFORM_FIELD_PLACEHOLDER}
              aria-label={CREATE_JOB_TEXTS.TRANSFORM_FIELD_ARIA}
              {...register(`transforms.${i}.params.field`)}
              className="flex-1"
            />
            <AppButton
              type="button"
              onClick={() => remove(i)}
              aria-label={CREATE_JOB_TEXTS.REMOVE_TRANSFORM_ARIA}
              className="px-3 text-neutral-400 hover:text-danger-600"
            >
              {COMMON_LABELS.REMOVE}
            </AppButton>
          </div>
          {errors.transforms?.[i]?.params?.field && (
            <p className="text-xs text-danger-600 mt-1">{errors.transforms[i]?.params?.field?.message}</p>
          )}
        </div>
      ))}
      <AppButton
        type="button"
        onClick={() => append({ name: 'uppercase', params: { field: '' } })}
        className="text-sm text-brand-600 hover:text-brand-700 font-medium"
      >
        {CREATE_JOB_TEXTS.ADD_TRANSFORM}
      </AppButton>
    </div>
  );
}
