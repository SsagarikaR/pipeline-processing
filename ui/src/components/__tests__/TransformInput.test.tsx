import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import TransformInput from '../TransformInput';

describe('TransformInput', () => {
  it('renders correctly and handles additions', () => {
    const onChange = vi.fn();
    render(<TransformInput transforms={[]} onChange={onChange} />);
    
    expect(screen.getByText('Transforms')).toBeInTheDocument();
    
    fireEvent.click(screen.getByText('+ Add transform'));
    expect(onChange).toHaveBeenCalledWith([{ name: 'uppercase', params: { field: '' } }]);
  });
});
