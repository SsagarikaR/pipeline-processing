import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import AggregationInput from '../AggregationInput';

describe('AggregationInput', () => {
  it('renders correctly and handles additions', () => {
    const onChange = vi.fn();
    render(<AggregationInput aggregations={[]} onChange={onChange} />);
    
    expect(screen.getByText('Aggregations')).toBeInTheDocument();
    
    fireEvent.click(screen.getByText('+ Add aggregation'));
    expect(onChange).toHaveBeenCalledWith([{ field: '', op: 'sum', groupBy: '', type: 'Number' }]);
  });
});
