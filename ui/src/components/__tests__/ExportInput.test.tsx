import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import ExportInput from '../ExportInput';

describe('ExportInput', () => {
  it('renders correctly and handles additions', () => {
    const onChange = vi.fn();
    render(<ExportInput exports={[]} onChange={onChange} />);
    
    expect(screen.getByText('Exports')).toBeInTheDocument();
    
    fireEvent.click(screen.getByText('+ Add export'));
    expect(onChange).toHaveBeenCalledWith([{ type: 's3', path: '' }]);
  });
});
