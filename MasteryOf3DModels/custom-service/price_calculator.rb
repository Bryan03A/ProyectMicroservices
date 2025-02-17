# price_calculator.rb

class PriceCalculator
  def self.calculate(price, custom_params)
    final_cost = price # The initial price is the model's price

    custom_params.each do |key, value|
      # Adjustment for filament quality
      if key == 'filament_quality'
        case value
        when 'Low'
          final_cost += 0
        when 'Medium'
          final_cost += 5.0
        when 'High'
          final_cost += 10.0
        end
      end

      # Adjustment for color
      if key == 'color'
        case value
        when 'Red'
          final_cost += 0
        when 'Gold'
          final_cost += 5.0
        when 'Blue'
          final_cost += 3.0
        end
      end

      # Adjustment for size
      if key == 'size'
        case value
        when 'Small'
          final_cost += 0
        when 'Medium'
          final_cost += 5.0
        when 'Large'
          final_cost += 8.0
        end
      end

      # Adjustment for shape complexity
      if key == 'shape_complexity'
        case value
        when 'Low'
          final_cost += 0
        when 'Medium'
          final_cost += 6.0
        when 'High'
          final_cost += 12.0
        end
      end

      # Adjustment for material type
      if key == 'material_type'
        case value
        when 'Standard'
          final_cost += 0
        when 'Premium'
          final_cost += 15.0
        when 'Luxury'
          final_cost += 25.0
        end
      end
    end

    return final_cost
  end
end