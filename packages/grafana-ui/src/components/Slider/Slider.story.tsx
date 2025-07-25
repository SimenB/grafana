import { StoryFn, Meta } from '@storybook/react';

import { Slider } from './Slider';

const meta: Meta<typeof Slider> = {
  title: 'Inputs/Slider',
  component: Slider,
  parameters: {
    controls: {
      exclude: ['formatTooltipResult', 'onChange', 'onAfterChange', 'value', 'tooltipAlwaysVisible'],
    },
    knobs: {
      disabled: true,
    },
    // TODO fix a11y issue in story and remove this
    a11y: { test: 'off' },
  },
  argTypes: {
    orientation: { control: { type: 'select', options: ['horizontal', 'vertical'] } },
    step: { control: { type: 'number', min: 1 } },
  },
  args: {
    min: 0,
    max: 100,
    value: 10,
    orientation: 'horizontal',
    reverse: false,
    included: true,
    step: undefined,
  },
};

export const Basic: StoryFn<typeof Slider> = (args) => {
  return (
    <div style={{ width: '300px', height: '300px' }}>
      <Slider {...args} />
    </div>
  );
};

export const WithMarks: StoryFn<typeof Slider> = (args) => {
  return (
    <div style={{ width: '300px', height: '300px' }}>
      <Slider {...args} />
    </div>
  );
};
WithMarks.args = {
  marks: { 0: '0', 25: '25', 50: '50', 75: '75', 100: '100' },
};

export default meta;
