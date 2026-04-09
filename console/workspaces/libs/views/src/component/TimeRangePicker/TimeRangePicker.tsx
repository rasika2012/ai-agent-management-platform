import {
  TraceListTimeRange,
  getTimeRange,
} from '@agent-management-platform/types';
import {
  AdapterDateFns,
  Box,
  DatePickers,
  Divider,
  Form,
  ListItemText,
  Menu,
  MenuItem,
  MenuList,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from '@wso2/oxygen-ui';
import { format } from 'date-fns';
import { ChevronDown, ChevronUp, Clock } from 'lucide-react';
import { useMemo, useState } from 'react';

export interface CustomTimeRange {
  startTime: string;
  endTime: string;
}
export type TimeRange = CustomTimeRange | TraceListTimeRange;

export function convertToTimeRange(timeRange: TimeRange): TraceListTimeRange {
  if (typeof timeRange === 'string') {
    return timeRange as TraceListTimeRange;
  }
  return timeRange;
}

const getTimeRangeLabel = (timeRange: TimeRange) => {
  if (typeof timeRange === 'string') {
    return timeRange;
  }
  return `${format(new Date(timeRange.startTime), 'dd/MM/yyyy, HH:mm')} - ${format(new Date(timeRange.endTime), 'dd/MM/yyyy, HH:mm')}`;
};
export const TIME_RANGE_OPTIONS = [
  { value: TraceListTimeRange.TEN_MINUTES, label: 'Last 10 Minutes' },
  { value: TraceListTimeRange.ONE_HOUR, label: 'Last 1 Hour' },
  { value: TraceListTimeRange.SIX_HOURS, label: 'Last 6 Hours' },
  { value: TraceListTimeRange.TWELVE_HOURS, label: 'Last 12 Hours' },
  { value: TraceListTimeRange.ONE_DAY, label: 'Last 1 Day' },
  { value: TraceListTimeRange.SEVEN_DAYS, label: 'Last 7 Days' },
];
export interface TimeRangePickerProps {
  timeRange: TimeRange;
  onChange: (timeRange: TimeRange) => void;
}
export function TimeRangePicker({ timeRange, onChange }: TimeRangePickerProps) {
  const normalizedRange = useMemo<CustomTimeRange>(() => {
    if (typeof timeRange === 'string') {
      return getTimeRange(timeRange as TraceListTimeRange);
    }
    return timeRange;
  }, [timeRange]);

  const startDate = useMemo(
    () => new Date(normalizedRange.startTime),
    [normalizedRange.startTime]
  );
  const endDate = useMemo(
    () => new Date(normalizedRange.endTime),
    [normalizedRange.endTime]
  );

  const emitCustomRange = (nextStart: Date | null, nextEnd: Date | null) => {
    if (!nextStart || !nextEnd || nextStart > nextEnd) return;
    onChange({
      startTime: nextStart.toISOString(),
      endTime: nextEnd.toISOString(),
    });
  };

  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const isOpen = Boolean(anchorEl);

  return (
    <Stack direction="row" spacing={1} alignItems="center">
   
   <Tooltip title={typeof timeRange === 'string' ? timeRange : getTimeRangeLabel(timeRange)}>
      <TextField
        value={
          typeof timeRange === 'string'
            ? timeRange
            : getTimeRangeLabel(timeRange)
        }
        onClick={(event) => {
          event.preventDefault();
          event.stopPropagation();
          setAnchorEl(event.currentTarget);
        }}
        fullWidth
        slotProps={{
          input: {
            readOnly: true,
            endAdornment: isOpen ? <ChevronUp size={28} /> : <ChevronDown size={28} />,
            startAdornment: <Box p={0.5} pr={1} alignItems='center' color="secondary" display='flex'><Clock size={18} /></Box>,
          },
        }}
      />
      </Tooltip>

      <Menu
        anchorEl={anchorEl}
        open={isOpen}
        onClose={() => setAnchorEl(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
        transformOrigin={{ vertical: 'top', horizontal: 'left' }}
      >
   
        <DatePickers.LocalizationProvider dateAdapter={AdapterDateFns}>
          <Stack spacing={1} sx={{ minWidth: 280 }}>

            <Stack spacing={1} px={2} py={1}>
            <Typography variant='h6' color='text.secondary'>Select Time Range</Typography>
              <Form.ElementWrapper name="startTime" label="Start Time">
                <DatePickers.DateTimePicker
                  value={startDate}
                  onChange={(value) => {
                    if (!value) return;
                    emitCustomRange(value, endDate);
                  }}
                  slotProps={{
                    textField: {
                      size: 'small',
                    },
                  }}
                />
              </Form.ElementWrapper>
              <Form.ElementWrapper name="endTime" label="End Time">
                <DatePickers.DateTimePicker
                  value={endDate}
                  onChange={(value) => {
                    if (!value) return;
                    emitCustomRange(startDate, value);
                  }}
                  slotProps={{
                    textField: {
                      size: 'small',
                    },
                  }}
                />
              </Form.ElementWrapper>
            </Stack>
            <Divider />
          </Stack>
          <MenuList>
            {TIME_RANGE_OPTIONS.map((option) => (
              <MenuItem
                key={option.value}
                selected={timeRange === option.value}
                onClick={() => {
                  onChange(option.value);
                  setAnchorEl(null);
                }}
              >
                <ListItemText>{option.label}</ListItemText>
              </MenuItem>
            ))}
          </MenuList>
        </DatePickers.LocalizationProvider>
      </Menu>
    </Stack>
  );
}
