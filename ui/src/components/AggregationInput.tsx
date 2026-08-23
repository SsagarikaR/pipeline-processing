import { useState } from 'react';
import { useFieldArray, useFormContext } from 'react-hook-form';
import type { AggregationOp } from '../types/job';
import type { JobSpecFormValues } from '../schemas/jobSpec';

const FIELD_TYPES = ['Number', 'String', 'Boolean', 'Date'] as const;

const AGG_OPS: Record<string, AggregationOp[]> = {
  Number: ['sum', 'avg', 'min', 'max', 'count'],
  String: ['count'],
  Boolean: ['count'],
  Date: ['min', 'max', 'count'],
};

export default function AggregationInput() {
  const {
    register,
    control,
    watch,
    setValue,
    formState: { errors },
  } = useFormContext<JobSpecFormValues>();
  const { fields, append, remove } = useFieldArray({ control, name: 'aggregations' });
  // UI-only grouping of ops by a coarse field type; not part of JobSpec, keyed by
  // the field array's stable id so it survives add/remove without reindexing.
  const [typeByFieldId, setTypeByFieldId] = useState<Record<string, string>>({});
  const currentOps = watch('aggregations');

  function handleTypeChange(fieldId: string, index: number, type: string) {
    setTypeByFieldId((prev) => ({ ...prev, [fieldId]: type }));
    const availableOps = AGG_OPS[type] ?? [];
    const currentOp = currentOps?.[index]?.op;
    if (!currentOp || !availableOps.includes(currentOp)) {
      setValue(`aggregations.${index}.op`, availableOps[0] ?? 'count', { shouldValidate: true });
    }
  }

  return (
    <div className="space-y-2">
      <label className="block text-sm font-medium text-neutral-700">
        Aggregations <span className="text-danger-600">*</span>
      </label>
      {fields.map((field, i) => {
        const currentType = typeByFieldId[field.id] ?? 'Number';
        const availableOps = AGG_OPS[currentType] ?? AGG_OPS.Number;
        const fieldError = errors.aggregations?.[i]?.field?.message;

        return (
          <div key={field.id}>
            <div className="flex gap-2">
              <select
                value={currentType}
                onChange={(e) => handleTypeChange(field.id, i, e.target.value)}
                className="border border-neutral-300 rounded-lg px-3 py-2 text-sm w-32"
              >
                {FIELD_TYPES.map((t) => (
                  <option key={t} value={t}>
                    {t}
                  </option>
                ))}
              </select>
              <select
                {...register(`aggregations.${i}.op`)}
                className="border border-neutral-300 rounded-lg px-3 py-2 text-sm w-28"
              >
                {availableOps.map((op) => (
                  <option key={op} value={op}>
                    {op}
                  </option>
                ))}
              </select>
              <input
                placeholder="field"
                {...register(`aggregations.${i}.field`)}
                className="border border-neutral-300 rounded-lg px-3 py-2 text-sm w-32"
              />
              <input
                placeholder="group by (optional)"
                {...register(`aggregations.${i}.groupBy`)}
                className="flex-1 border border-neutral-300 rounded-lg px-3 py-2 text-sm"
              />
              <button type="button" onClick={() => remove(i)} className="px-3 text-neutral-400 hover:text-danger-600">
                ✕
              </button>
            </div>
            {fieldError && <p className="text-xs text-danger-600 mt-1">{fieldError}</p>}
          </div>
        );
      })}
      {errors.aggregations?.message && <p className="text-xs text-danger-600">{errors.aggregations.message}</p>}
      <button
        type="button"
        onClick={() => append({ field: '', op: 'sum', groupBy: '' })}
        className="text-sm text-brand-600 hover:text-brand-700 font-medium"
      >
        + Add aggregation
      </button>
    </div>
  );
}
