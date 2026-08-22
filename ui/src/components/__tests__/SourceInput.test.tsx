import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import SourceInput from '../SourceInput';

describe('SourceInput', () => {
  it('renders correctly and handles additions', () => {
    const onChange = vi.fn();
    render(<SourceInput sources={[]} onChange={onChange} />);
    
    expect(screen.getByText('Sources')).toBeInTheDocument();
    
    fireEvent.click(screen.getByText('+ Add source'));
    expect(onChange).toHaveBeenCalledWith([{ type: 'csv', path: '' }]);
  });
});
