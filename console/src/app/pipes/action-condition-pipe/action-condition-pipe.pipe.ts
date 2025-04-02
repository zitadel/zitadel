import { Pipe, PipeTransform } from '@angular/core';
import { Condition } from '@zitadel/proto/zitadel/action/v2beta/execution_pb';

@Pipe({
  name: 'condition',
})
export class ActionConditionPipe implements PipeTransform {
  transform(condition?: Condition): string {
    if (!condition?.conditionType?.case) {
      return '';
    }

    const conditionType = condition.conditionType.value;

    if ('name' in conditionType) {
      // Applies for function condition
      return `function: ${conditionType.name}`;
    }

    const { condition: innerCondition } = conditionType;

    if (typeof innerCondition.value === 'string') {
      // Applies for service, method condition of Request/ResponseCondition, event, and group of EventCondition
      return `${innerCondition.case}: ${innerCondition.value}`;
    }

    return `all`;
  }
}
